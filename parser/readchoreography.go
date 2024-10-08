package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/choreography"
	"github.com/ulikunitz/xz/lzma"
)

const BVCD_MAGIC = 0x64637662 // bvcd

func parseVcdList(context *parseContext, block *source2.FileBlock) error {
	vcdList := &source2.FileBlockVcdList{}
	var choreoCount uint32
	var err error
	var strings []string
	var o int64

	if err = binary.Read(context.reader, binary.LittleEndian, &vcdList.Version); err != nil {
		return fmt.Errorf("failed to read vcd list version: <%w>", err)
	}
	if err = binary.Read(context.reader, binary.LittleEndian, &choreoCount); err != nil {
		return fmt.Errorf("failed to read vcd list count: <%w>", err)
	}

	o, _ = context.reader.Seek(0, io.SeekCurrent)
	if strings, err = readVcdStrings(context.reader); err != nil {
		return err
	}

	context.reader.Seek(o+16, io.SeekStart) //skip strings data
	if vcdList.Choreographies, err = readChoreographies(context.reader, choreoCount, strings); err != nil {
		return err
	}

	block.Content = vcdList

	return nil
}

func readVcdStrings(reader io.ReadSeeker) ([]string, error) {
	var err error
	var offset uint32
	var stringsCount uint32
	var stringsOffset uint32
	var o int64

	o, _ = reader.Seek(0, io.SeekCurrent)
	if err = binary.Read(reader, binary.LittleEndian, &offset); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings offset offset: <%w>", err)
	}
	offset += uint32(o)
	if err = binary.Read(reader, binary.LittleEndian, &stringsCount); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings count: <%w>", err)
	}
	o, _ = reader.Seek(0, io.SeekCurrent)
	if err = binary.Read(reader, binary.LittleEndian, &stringsOffset); err != nil {
		return nil, fmt.Errorf("failed to read vcd strings offset: <%w>", err)
	}
	stringsOffset += uint32(o)

	stringOffsets := make([]uint32, stringsCount)
	strings := make([]string, stringsCount)

	reader.Seek(int64(offset), io.SeekStart)
	for i := uint32(0); i < stringsCount; i++ {
		if err = binary.Read(reader, binary.LittleEndian, &stringOffsets[i]); err != nil {
			return nil, fmt.Errorf("failed to read string offset %d: <%w>", i, err)
		}
	}

	for i := uint32(0); i < stringsCount; i++ {
		reader.Seek(int64(stringsOffset+stringOffsets[i]), io.SeekStart)

		if strings[i], err = readNullTerminatedString(reader); err != nil {
			return nil, fmt.Errorf("failed to read string %d: <%w>", i, err)
		}
	}
	return strings, nil
}

func readChoreographies(reader io.ReadSeeker, count uint32, strings []string) ([]*choreography.Choreography, error) {
	var err error
	var o int64
	choreos := make([]*choreography.Choreography, count)
	const CHOREO_LENGTH = 24
	o, _ = reader.Seek(0, io.SeekCurrent)
	for i := uint32(0); i < count; i++ {
		reader.Seek(o+int64(i*CHOREO_LENGTH), io.SeekStart)
		if choreos[i], err = readChoreography(reader, strings); err != nil {
			return nil, err
		}
	}

	return choreos, nil
}

func readChoreography(reader io.ReadSeeker, strings []string) (*choreography.Choreography, error) {
	choreo := choreography.NewChoreography()
	var nameOffset uint32
	var blockOffset uint32
	var blockLength uint32
	var hasSounds uint32
	var blockReader io.ReadSeeker
	var err error

	if nameOffset, err = readOffset(reader); err != nil {
		return nil, fmt.Errorf("failed to read name offset: <%w>", err)
	}
	if blockOffset, err = readOffset(reader); err != nil {
		return nil, fmt.Errorf("failed to read block offset: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &blockLength); err != nil {
		return nil, fmt.Errorf("failed to read choreography length: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &choreo.Duration); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &choreo.SoundDuration); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &hasSounds); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	choreo.HasSounds = hasSounds == 1

	reader.Seek(int64(nameOffset), io.SeekStart)
	if choreo.Name, err = readNullTerminatedString(reader); err != nil {
		return nil, fmt.Errorf("failed to read choreography name: <%w>", err)
	}

	reader.Seek(int64(blockOffset), io.SeekStart)
	if blockReader, err = readBlock(reader, blockLength); err != nil {
		return nil, fmt.Errorf("failed to read choreography block: <%w>", err)
	}

	if err = readChoreographyData(blockReader, strings, choreo); err != nil {
		return nil, err
	}

	//log.Println(blockReader)
	//log.Println(nameOffset, blockOffset)

	return choreo, nil
}

func readChoreographyData(reader io.ReadSeeker, strings []string, choreo *choreography.Choreography) error {
	var err error
	var bvcdMagic uint32
	var version uint8
	if err = binary.Read(reader, binary.LittleEndian, &bvcdMagic); err != nil {
		return fmt.Errorf("failed to read choreography bvcd magic: <%w>", err)
	}

	if bvcdMagic != BVCD_MAGIC {
		return fmt.Errorf("unsupported choreography magic: %d", bvcdMagic)
	}

	if err = binary.Read(reader, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("failed to read choreography version: <%w>", err)
	}
	if version > 19 {
		return fmt.Errorf("unsupported choreography version: <%d>", version)
	}

	if _, err = reader.Seek(4, io.SeekCurrent); err != nil { // Skip crc
		return fmt.Errorf("failed to skip crc: <%w>", err)
	}

	if err = readChoreographyEvents(reader, strings, choreo, version); err != nil {
		return fmt.Errorf("failed to read choreography events: <%w>", err)
	}

	if err = readChoreographyActors(reader, strings, choreo, version); err != nil {
		return fmt.Errorf("failed to read choreography actors: <%w>", err)
	}

	if choreo.CurveData, err = readCurveData(reader, version); err != nil {
		return fmt.Errorf("failed to read choreography curve data: <%w>", err)
	}

	var ignorePhonemes uint8
	if err = binary.Read(reader, binary.LittleEndian, &ignorePhonemes); err != nil {
		return fmt.Errorf("failed to read choreography ignore phonemes: <%w>", err)
	}
	choreo.IgnorePhonemes = ignorePhonemes == 1

	return nil
}

func readChoreographyEvents(reader io.ReadSeeker, strings []string, choreo *choreography.Choreography, version uint8) error {
	var eventCount uint8
	var err error
	var event *choreography.ChoreographyEvent
	if err = binary.Read(reader, binary.LittleEndian, &eventCount); err != nil {
		return fmt.Errorf("failed to read choreography event count: <%w>", err)
	}

	for i := 0; i < int(eventCount); i++ {
		if event, err = readChoreographyEvent(reader, strings, version); err != nil {
			return fmt.Errorf("failed to read choreography event: <%w>", err)
		}
		choreo.AddEvent(event)
	}

	return nil
}

func readChoreographyEvent(reader io.ReadSeeker, strings []string, version uint8) (*choreography.ChoreographyEvent, error) {
	event := choreography.NewChoreographyEvent()
	var err error
	var count uint8
	var tag *choreography.ChoreographyTag

	if err = binary.Read(reader, binary.LittleEndian, &event.EventType); err != nil {
		return nil, fmt.Errorf("failed to read event type: <%w>", err)
	}
	if event.Name, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read event name: <%w>", err)
	}
	if event.PreferredName, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read event name: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &event.StartTime); err != nil {
		return nil, fmt.Errorf("failed to read event start time: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &event.EndTime); err != nil {
		return nil, fmt.Errorf("failed to read event end time: <%w>", err)
	}
	if event.Param1, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read event param 1: <%w>", err)
	}
	if event.Param2, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read event param 2: <%w>", err)
	}
	if event.Param3, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read event param 3: <%w>", err)
	}
	if event.CurveData, err = readCurveData(reader, version); err != nil {
		return nil, fmt.Errorf("failed to read event curve data: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &event.Flags); err != nil {
		return nil, fmt.Errorf("failed to read event flags: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &event.DistToTarget); err != nil {
		return nil, fmt.Errorf("failed to read event dist to target: <%w>", err)
	}

	// Relative tag
	if err = binary.Read(reader, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("failed to read event tag count: <%w>", err)
	}
	event.RelativeTags = make([]*choreography.ChoreographyTag, count)
	for i := 0; i < int(count); i++ {
		if tag, err = readChoreographyTag(reader, strings); err != nil {
			return nil, fmt.Errorf("failed to read choreography tag: <%w>", err)
		}
		event.RelativeTags[i] = tag
	}

	// Flex timing tag
	if err = binary.Read(reader, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("failed to read event tag count: <%w>", err)
	}
	event.FlexTimingTags = make([]*choreography.ChoreographyTag, count)
	for i := 0; i < int(count); i++ {
		if tag, err = readChoreographyTag(reader, strings); err != nil {
			return nil, fmt.Errorf("failed to read choreography tag: <%w>", err)
		}
		event.FlexTimingTags[i] = tag
	}

	// absolute tags
	for j := 0; j < 2; j++ {
		if err = binary.Read(reader, binary.LittleEndian, &count); err != nil {
			return nil, fmt.Errorf("failed to read event tag count: <%w>", err)
		}
		event.AbsoluteTags[j] = make([]*choreography.ChoreographyTag, count)
		for i := 0; i < int(count); i++ {
			if tag, err = readChoreographyAbsoluteTag(reader, strings); err != nil {
				return nil, fmt.Errorf("failed to read choreography tag: <%w>", err)
			}
			event.AbsoluteTags[j][i] = tag
		}
	}

	if event.EventType == choreography.GESTURE {
		var duration float32

		if err = binary.Read(reader, binary.LittleEndian, &duration); err != nil {
			return nil, fmt.Errorf("failed to read event duration: <%w>", err)
		}
	}

	var useTag uint8
	if err = binary.Read(reader, binary.LittleEndian, &useTag); err != nil {
		return nil, fmt.Errorf("failed to read choreography duration: <%w>", err)
	}
	if useTag == 1 {
		event.UsesTag = true
		if event.TagName, err = readString(reader, strings); err != nil {
			return nil, fmt.Errorf("failed to read event tag name: <%w>", err)
		}
		if event.TagWavName, err = readString(reader, strings); err != nil {
			return nil, fmt.Errorf("failed to read event tag wav name: <%w>", err)
		}
	}

	if err = readChoreographyFlexAnimations(reader, strings, event, version); err != nil {
		return nil, err
	}

	if event.EventType == choreography.LOOP {
		var loopCount int8
		if err = binary.Read(reader, binary.LittleEndian, &loopCount); err != nil {
			return nil, fmt.Errorf("failed to read event loop count: <%w>", err)
		}
		event.LoopCount = max(-1, int(loopCount))
	}

	if event.EventType == choreography.SPEAK {
		if version < 17 {
			if err = binary.Read(reader, binary.LittleEndian, &event.CloseCaptionType); err != nil {
				return nil, fmt.Errorf("failed to read event close caption type: <%w>", err)
			}
			if event.CloseCaptionToken, err = readString(reader, strings); err != nil {
				return nil, fmt.Errorf("failed to read event close caption token: <%w>", err)
			}
		}
		if err = binary.Read(reader, binary.LittleEndian, &event.SpeakFlags); err != nil {
			return nil, fmt.Errorf("failed to read event speak flags: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &event.SoundStartDelay); err != nil {
			return nil, fmt.Errorf("failed to read event sound start delay: <%w>", err)
		}
	}

	if err = binary.Read(reader, binary.LittleEndian, &event.ConstrainedEventId); err != nil {
		return nil, fmt.Errorf("failed to read event constrained event id: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &event.EventId); err != nil {
		return nil, fmt.Errorf("failed to read event event id: <%w>", err)
	}

	return event, nil
}

func readChoreographyTag(reader io.ReadSeeker, strings []string) (*choreography.ChoreographyTag, error) {
	var err error
	var value uint8
	tag := &choreography.ChoreographyTag{}

	if tag.Name, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read tag name: <%w>", err)
	}

	if err = binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return nil, fmt.Errorf("failed to read event end time: <%w>", err)
	}

	tag.Value = float32(value) / 255.
	return tag, nil
}

func readChoreographyAbsoluteTag(reader io.ReadSeeker, strings []string) (*choreography.ChoreographyTag, error) {
	var err error
	var value uint16
	tag := &choreography.ChoreographyTag{}

	if tag.Name, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read tag name: <%w>", err)
	}

	if err = binary.Read(reader, binary.LittleEndian, &value); err != nil {
		return nil, fmt.Errorf("failed to read event end time: <%w>", err)
	}

	tag.Value = float32(value) / 4096.
	return tag, nil
}

func readChoreographyFlexAnimations(reader io.ReadSeeker, strings []string, event *choreography.ChoreographyEvent, version uint8) error {
	var numTracks byte
	var controllerName string
	var track *choreography.FlexAnimationTrack
	var err error

	if err = binary.Read(reader, binary.LittleEndian, &numTracks); err != nil {
		return fmt.Errorf("failed to read event track count: <%w>", err)
	}

	for i := 0; i < int(numTracks); i++ {
		if controllerName, err = readString(reader, strings); err != nil {
			return fmt.Errorf("failed to read event controller name: <%w>", err)
		}
		track = event.AddTrack(controllerName)

		if err = binary.Read(reader, binary.LittleEndian, &track.Flags); err != nil {
			return fmt.Errorf("failed to read event track flags: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &track.Min); err != nil {
			return fmt.Errorf("failed to read event track min: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &track.Max); err != nil {
			return fmt.Errorf("failed to read event track max: <%w>", err)
		}

		if track.SamplesCurve, err = readCurveData(reader, version); err != nil {
			return fmt.Errorf("failed to read event curve data: <%w>", err)
		}

		if track.IsComboType() {
			if track.ComboSamplesCurve, err = readCurveData(reader, version); err != nil {
				return fmt.Errorf("failed to read event curve data: <%w>", err)
			}
		}
	}

	return nil
}

func readString(reader io.ReadSeeker, strings []string) (string, error) {
	var index uint32
	if err := binary.Read(reader, binary.LittleEndian, &index); err != nil {
		return "", fmt.Errorf("failed to read string index: <%w>", err)
	}

	if int(index) >= len(strings) {
		return "", fmt.Errorf("index out of bounds: %d >=  %d", index, len(strings))
	}

	return strings[index], nil
}

func readCurveData(reader io.ReadSeeker, version uint8) (*choreography.CurveData, error) {
	curveData := choreography.NewCurveData()
	var count uint16
	var err error
	var sample *choreography.CurveDataSample
	var sampleType uint8
	var unk uint8

	if err = binary.Read(reader, binary.LittleEndian, &count); err != nil {
		return nil, fmt.Errorf("failed to read curve data count: <%w>", err)
	}

	for i := 0; i < int(count); i++ {
		if sampleType == 0 {
			if sample, err = readCurveDataSample(reader); err != nil {
				return nil, fmt.Errorf("failed to read curve data time: <%w>", err)
			}
			curveData.AddSample(sample)
		} else if sampleType == 1 {
			curveType := &choreography.CurveDataSampleType{}
			if err = binary.Read(reader, binary.LittleEndian, &curveType.OutType); err != nil {
				return nil, fmt.Errorf("failed to read curve data out type: <%w>", err)
			}
			if err = binary.Read(reader, binary.LittleEndian, &curveType.InType); err != nil {
				return nil, fmt.Errorf("failed to read curve data in type: <%w>", err)
			}
			if err = binary.Read(reader, binary.LittleEndian, &unk); err != nil {
				return nil, fmt.Errorf("failed to read curve data unk: <%w>", err)
			}

			sample.CurveType = curveType

			/*

			   var outType = reader.ReadByte();
			   var inType = reader.ReadByte();
			   var nullTermination = reader.ReadByte();
			   Debug.Assert(nullTermination == 0);

			   lastSample.SetCurveType(inType, outType);
			*/

		} else {
			return nil, errors.New("unknonw sample type")
		}
		if err = binary.Read(reader, binary.LittleEndian, &sampleType); err != nil {
			return nil, fmt.Errorf("failed to read curve data sample type: <%w>", err)
		}
	}

	if version >= 16 {
		if curveData.LeftEdge, err = readCurveDataEdge(reader); err != nil {
			return nil, fmt.Errorf("failed to read left edge: <%w>", err)
		}
		if curveData.RightEdge, err = readCurveDataEdge(reader); err != nil {
			return nil, fmt.Errorf("failed to read right edge: <%w>", err)
		}
	}

	return curveData, nil
}

func readCurveDataSample(reader io.ReadSeeker) (*choreography.CurveDataSample, error) {
	sample := &choreography.CurveDataSample{}
	var err error
	var t float32
	var v uint8
	var hasBezier uint8

	if err = binary.Read(reader, binary.LittleEndian, &t); err != nil {
		return nil, fmt.Errorf("failed to read curve data time: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &v); err != nil {
		return nil, fmt.Errorf("failed to read curve data value: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &hasBezier); err != nil {
		return nil, fmt.Errorf("failed to read curve data has bezier: <%w>", err)
	}

	sample.Time = t
	sample.Value = float32(v) / 255.0

	if hasBezier == 1 {
		sample.Bezier = &choreography.CurveDataSampleBezier{}
		if err = binary.Read(reader, binary.LittleEndian, &sample.Bezier.Flags); err != nil {
			return nil, fmt.Errorf("failed to read curve data bezier flags: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &sample.Bezier.InDeg); err != nil {
			return nil, fmt.Errorf("failed to read curve data bezier in: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &sample.Bezier.InWeight); err != nil {
			return nil, fmt.Errorf("failed to read curve data bezier in weight: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &sample.Bezier.OutDeg); err != nil {
			return nil, fmt.Errorf("failed to read curve data bezier out: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &sample.Bezier.OutWeight); err != nil {
			return nil, fmt.Errorf("failed to read curve data bezier out weight: <%w>", err)
		}
	}

	return sample, nil
}

func readCurveDataEdge(reader io.ReadSeeker) (*choreography.CurveDataEdge, error) {
	var hasEdge uint8
	var err error

	if err = binary.Read(reader, binary.LittleEndian, &hasEdge); err != nil {
		return nil, fmt.Errorf("failed to read curve data has edge: <%w>", err)
	}
	if hasEdge == 0 {
		return nil, nil
	}

	edge := &choreography.CurveDataEdge{}
	if err = binary.Read(reader, binary.LittleEndian, &edge.CurveType.InType); err != nil {
		return nil, fmt.Errorf("failed to read choreography edge in type: <%w>", err)
	}
	if err = binary.Read(reader, binary.LittleEndian, &edge.CurveType.OutType); err != nil {
		return nil, fmt.Errorf("failed to read choreography edge out type: <%w>", err)
	}
	reader.Seek(2, io.SeekCurrent)
	if err = binary.Read(reader, binary.LittleEndian, &edge.ZeroValue); err != nil {
		return nil, fmt.Errorf("failed to read choreography edge zero value: <%w>", err)
	}

	return edge, nil
}

func readChoreographyActors(reader io.ReadSeeker, strings []string, choreo *choreography.Choreography, version uint8) error {
	var actorCount uint8
	var err error
	var actor *choreography.ChoreographyActor
	if err = binary.Read(reader, binary.LittleEndian, &actorCount); err != nil {
		return fmt.Errorf("failed to read choreography actor count: <%w>", err)
	}

	for i := 0; i < int(actorCount); i++ {
		if actor, err = readChoreographyActor(reader, strings, version); err != nil {
			return fmt.Errorf("failed to read choreography actor: <%w>", err)
		}
		choreo.AddActor(actor)
	}

	return nil
}

func readChoreographyActor(reader io.ReadSeeker, strings []string, version uint8) (*choreography.ChoreographyActor, error) {
	actor := choreography.NewChoreographyActor()
	//channel := choreography.NewChoreographyChannel()
	var err error
	var channelCount uint8
	var active uint8
	var channel *choreography.ChoreographyChannel

	if actor.Name, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read actor name: <%w>", err)
	}

	if err = binary.Read(reader, binary.LittleEndian, &channelCount); err != nil {
		return nil, fmt.Errorf("failed to read choreography actor channel count: <%w>", err)
	}

	for i := 0; i < int(channelCount); i++ {
		if channel, err = readChoreographyChannel(reader, strings, version); err != nil {
			return nil, fmt.Errorf("failed to read choreography actor: <%w>", err)
		}
		actor.AddChannel(channel)
	}

	if err = binary.Read(reader, binary.LittleEndian, &active); err != nil {
		return nil, fmt.Errorf("failed to read choreography actor active: <%w>", err)
	}
	actor.Active = active == 1

	return actor, nil
}

func readChoreographyChannel(reader io.ReadSeeker, strings []string, version uint8) (*choreography.ChoreographyChannel, error) {
	channel := choreography.NewChoreographyChannel()
	var err error
	var eventCount uint8
	var active uint8
	var event *choreography.ChoreographyEvent

	if channel.Name, err = readString(reader, strings); err != nil {
		return nil, fmt.Errorf("failed to read channel name: <%w>", err)
	}

	if err = binary.Read(reader, binary.LittleEndian, &eventCount); err != nil {
		return nil, fmt.Errorf("failed to read choreography channel event count: <%w>", err)
	}

	for i := 0; i < int(eventCount); i++ {
		if event, err = readChoreographyEvent(reader, strings, version); err != nil {
			return nil, fmt.Errorf("failed to read choreography actor: <%w>", err)
		}
		channel.AddEvent(event)
	}

	if err = binary.Read(reader, binary.LittleEndian, &active); err != nil {
		return nil, fmt.Errorf("failed to read choreography actor active: <%w>", err)
	}
	channel.Active = active == 1

	return channel, nil
}

func readOffset(reader io.ReadSeeker) (uint32, error) {
	var offset uint32
	o, _ := reader.Seek(0, io.SeekCurrent)
	if err := binary.Read(reader, binary.LittleEndian, &offset); err != nil {
		return 0, fmt.Errorf("failed to read offset: <%w>", err)
	}
	offset += uint32(o)
	return offset, nil
}

func readBlock(reader io.ReadSeeker, length uint32) (io.ReadSeeker, error) {
	var lzmaMagic uint32
	var err error
	var r *lzma.Reader

	if err = binary.Read(reader, binary.LittleEndian, &lzmaMagic); err != nil {
		return nil, fmt.Errorf("failed to read choreography block magic: <%w>", err)
	}

	if lzmaMagic == 0x414D5A4C { //LZMA
		var uncompressedSize uint32
		var compressedSize uint32
		//var properties [5]byte

		if err = binary.Read(reader, binary.LittleEndian, &uncompressedSize); err != nil {
			return nil, fmt.Errorf("failed to read lzma uncompressedSize: <%w>", err)
		}
		if err = binary.Read(reader, binary.LittleEndian, &compressedSize); err != nil {
			return nil, fmt.Errorf("failed to read lzma compressedSize: <%w>", err)
		}

		b := make([]byte, compressedSize+13)
		buf := bytes.NewBuffer(make([]byte, 0, 8))
		if err = binary.Read(reader, binary.LittleEndian, b[0:5]); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}
		if err = binary.Write(buf, binary.LittleEndian, uint64(uncompressedSize)); err != nil {
			return nil, fmt.Errorf("failed to write choreography uncompressed size: <%w>", err)
		}
		copy(b[5:], buf.Bytes())
		if err = binary.Read(reader, binary.LittleEndian, b[13:]); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}

		if r, err = lzma.NewReader(bytes.NewReader(b)); err != nil {
			return nil, fmt.Errorf("failed to create lzma reader: <%w>", err)
		}

		var uncompressed = make([]byte, uncompressedSize)
		if _, err = r.Read(uncompressed[:]); err != nil {
			return nil, err
		}

		return bytes.NewReader(uncompressed), nil
	} else {
		reader.Seek(-4, io.SeekCurrent)
		//TODO: optimize
		b := make([]byte, length)
		if err = binary.Read(reader, binary.LittleEndian, &b); err != nil {
			return nil, fmt.Errorf("failed to read choreography bytes: <%w>", err)
		}
		return bytes.NewReader(b), nil
	}
}

func readNullTerminatedString(reader io.ReadSeeker) (string, error) {
	builder := strings.Builder{}
	var b byte

	for {
		if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
			return "", fmt.Errorf("unable to read null string: <%w>", err)
		}
		if b == 0 {
			break
		}
		builder.WriteString(string(b))
	}

	return builder.String(), nil
}

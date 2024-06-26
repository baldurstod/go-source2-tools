package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/baldurstod/go-source2-tools/kv3"
	"github.com/pierrec/lz4/v4"
)

type parseKv3Context struct {
	reader           io.ReadSeeker
	stringDictionary []string
	decompressOffset uint32
}

func newParseKv3Context(reader io.ReadSeeker) *parseKv3Context {
	return &parseKv3Context{
		reader: reader,
	}
}

func ParseKv3(b []byte, version int, singleByteCount uint32, quadByteCount uint32, eightByteCount uint32, dictionaryTypeLength uint32,
	blobCount uint32, totalUncompressedBlobSize uint32, compressedBlobReader io.ReadSeeker, uncompressedBlobReader io.ReadSeeker, compressionFrameSize uint16) (*kv3.Kv3Element, error) {
	context := newParseKv3Context(bytes.NewReader(b))
	log.Println("Start parsing kv3", blobCount)

	quadCursor := math.Ceil(float64(singleByteCount)/4) * 4
	eightCursor := math.Ceil((quadCursor+float64(quadByteCount)*4)/8) * 8

	dictionaryOffset := uint32(eightCursor) + eightByteCount*8
	blobOffset := dictionaryOffset + dictionaryTypeLength
	blobEnd := uint32(len(b) - 4)

	var uncompressedBlobSizeReader, compressedBlobSizeReader io.ReadSeeker
	if version >= 2 && blobCount != 0 {
		if compressedBlobReader != nil {
			uncompressedLength := blobCount * 4
			compressedEnd := blobOffset + uncompressedLength
			uncompressedBlobSizeReader = bytes.NewReader(b[blobOffset : blobOffset+uncompressedLength])
			compressedBlobSizeReader = bytes.NewReader(b[compressedEnd+4 : compressedEnd+4+blobCount*2])

			log.Println(uncompressedBlobSizeReader, compressedBlobSizeReader)

		} else {
			if uncompressedBlobReader != nil {
				uncompressedBlobSizeReader = bytes.NewReader(b[blobEnd-blobCount*4 : blobEnd])
				log.Println(uncompressedBlobSizeReader)
			}
		}
	}

	var offset uint32
	if version == 1 {
		offset = blobEnd
	} else if version >= 2 {
		offset = blobOffset
	}

	// First compute the size
	size := 0
	s := offset
	for {
		s--
		if s <= 0 {
			break
		}
		if b[s] != 0 {
			size++
		} else {
			break
		}
	}

	typeReader := bytes.NewReader(b[s+1 : offset])
	valueArray := make([]kv3.Kv3Value, s)
	//let valueArray = [];
	//byteReader := bytes.NewReader(b)
	quadReader := bytes.NewReader(b)
	eightReader := bytes.NewReader(b)

	quadReader.Seek(int64(quadCursor), io.SeekStart)
	eightReader.Seek(int64(eightCursor), io.SeekStart)

	var stringCount uint32
	err := binary.Read(quadReader, binary.LittleEndian, &stringCount)
	if err != nil {
		return nil, fmt.Errorf("failed to read stringCount in ParseKv3: <%w>", err)
	}

	context.reader.Seek(int64(dictionaryOffset), io.SeekStart)
	context.stringDictionary = make([]string, stringCount)
	readStringDictionary(context, stringCount)

	var decompressBlobBuffer []byte
	//let decompressBlobArray;

	if compressedBlobReader != nil { //if a compressed reader is provided, we have to uncompress the blobs
		decompressBlobBuffer = make([]byte, totalUncompressedBlobSize) //new ArrayBuffer(totalUncompressedBlobSize);
		/*decompressBlobArray = new Uint8Array(decompressBlobBuffer);
		decompressBlobArray.decompressOffset = 0;*/
		context.decompressOffset = uint32(0)
	}

	elementType, err := readElementType(typeReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read elementType in ParseKv3: <%w>", err)
	}

	rootElement, err := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, compressedBlobReader, uncompressedBlobReader, typeReader, valueArray, elementType, false, compressionFrameSize)
	if err != nil {
		return nil, fmt.Errorf("call to parsebinarykv3element returned an error in ParseKv3: <%w>", err)
	}

	log.Println("End parsing kv3", size, stringCount, rootElement)
	return rootElement.(*kv3.Kv3Element), nil
}

func readStringDictionary(context *parseKv3Context, stringCount uint32) error {
	for i := uint32(0); i < stringCount; i++ {
		s, err := readNullString(context.reader)
		if err != nil {
			return fmt.Errorf("failed to read string in readStringDictionary: <%w>", err)
		}
		context.stringDictionary = append(context.stringDictionary, s)
	}
	log.Println(context.stringDictionary)
	return nil
}

func parseBinaryKv3Element(context *parseKv3Context, quadReader *bytes.Reader, eightReader *bytes.Reader, uncompressedBlobSizeReader io.ReadSeeker, compressedBlobSizeReader io.ReadSeeker, blobCount uint32,
	decompressBlobBuffer []byte, compressedBlobReader io.ReadSeeker, uncompressedBlobReader io.ReadSeeker, typeReader io.ReadSeeker, valueArray []kv3.Kv3Value, elementType byte, isArray bool, compressionFrameSize uint16) (kv3.Kv3Value, error) {

	switch elementType {
	case kv3.DATA_TYPE_NULL:
		return nil, nil
	case kv3.DATA_TYPE_BOOL:
		var b uint8
		err := binary.Read(context.reader, binary.LittleEndian, &b)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			//return byteReader.getUint8() ? true : false;
			if b > 0 {
				return true, nil
			} else {
				return false, nil
			}
		} else {
			valueArray = append(valueArray, b)
			/*
				let value = new SourceKv3Value(elementType);
				valueArray.push(value);
				value.value = byteReader.getUint8() ? true : false;
				return value;*/
			return nil, nil
		}
	case kv3.DATA_TYPE_INT64:
		var value int64
		err := binary.Read(eightReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			return value, nil
		} else {
			valueArray = append(valueArray, value)
			return nil, nil
		}
	case kv3.DATA_TYPE_UINT64:
		var value uint64
		err := binary.Read(eightReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			return value, nil
		} else {
			valueArray = append(valueArray, value)
			return nil, nil
		}
	case kv3.DATA_TYPE_DOUBLE:
		var value float64
		err := binary.Read(eightReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			return value, nil
		} else {
			valueArray = append(valueArray, value)
			return nil, nil
		}
	case kv3.DATA_TYPE_BYTE:
		var value int8
		err := binary.Read(context.reader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			return value, nil
		} else {
			valueArray = append(valueArray, value)
			return nil, nil
		}
	case kv3.DATA_TYPE_STRING, kv3.DATA_TYPE_RESOURCE:
		var value int32
		err := binary.Read(quadReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}

		return value, nil
	case kv3.DATA_TYPE_BLOB:
		if blobCount == 0 {
			var count uint32
			err := binary.Read(context.reader, binary.LittleEndian, &count)
			if err != nil {
				return nil, fmt.Errorf("failed to read count in parseBinaryKv3Element: <%w>", err)
			}
			//count = quadReader.getUint32();
			//elements = [];

			elements := make([]byte, count)
			var value byte

			// TODO: copy byte array as a whole
			log.Println(count)
			for i := uint32(0); i < count; i++ {
				err := binary.Read(context.reader, binary.LittleEndian, &value)
				if err != nil {
					return nil, fmt.Errorf("failed to read blob value in parseBinaryKv3Element: <%w>", err)
				}
				elements[i] = value
			}
			/*for (let i = 0; i < count; i++) {
				elements.push(byteReader.getUint8());
			}*/
			//return elements
			return nil, nil
		} else {
			if compressedBlobReader != nil { //if we have a decompress buffer, that means we have to decompress the blobs
				var uncompressedBlobSize uint32
				var compressedBlobSize uint16

				err := binary.Read(uncompressedBlobSizeReader, binary.LittleEndian, &uncompressedBlobSize)
				if err != nil {
					return nil, fmt.Errorf("failed to read uncompressedBlobSize in parseBinaryKv3Element: <%w>", err)
				}

				err = binary.Read(compressedBlobSizeReader, binary.LittleEndian, &compressedBlobSize)
				if err != nil {
					return nil, fmt.Errorf("failed to read compressedBlobSize in parseBinaryKv3Element: <%w>", err)
				}

				//let decompressBuffer = new ArrayBuffer(uncompressedBlobSize);

				end := context.decompressOffset + uncompressedBlobSize
				log.Println("decompressOffset", context.decompressOffset, end, compressionFrameSize, compressedBlobSize, uncompressedBlobSize)
				var decompressArray = decompressBlobBuffer[context.decompressOffset:end]
				var suboffset = uint32(0)
				for {
					src := make([]byte, compressedBlobSize)
					_, err := compressedBlobReader.Read(src)
					if err != nil {
						return nil, fmt.Errorf("failed to read block lz4 source in parseBinaryKv3Element: <%w>", err)
					}

					log.Println("suboffset", suboffset, uncompressedBlobSize, compressionFrameSize)
					var decompressArray2 = decompressArray[suboffset:uncompressedBlobSize]
					dict := decompressBlobBuffer[:context.decompressOffset]
					log.Println("decompressArray2", suboffset, uncompressedBlobSize, len(decompressArray2))
					if uncompressedBlobSize > uint32(compressionFrameSize) {
						uncompressedFrameSize, err := lz4.UncompressBlockWithDict(src, decompressArray2, dict)

						log.Println(src)
						if err != nil {
							return nil, fmt.Errorf("failed to decode lz4 <1> in parseBinaryKv3Element: <%w>", err)
						}
						//const uncompressedFrameSize = decodeLz4(compressedBlobReader, decompressBlobArray, compressedBlobSize, compressionFrameSize, decompressBlobArray.decompressOffset)
						context.decompressOffset += uint32(uncompressedFrameSize)
						suboffset += uint32(uncompressedFrameSize)
						uncompressedBlobSize -= uint32(uncompressedFrameSize)
					} else {
						//uncompressedBlobSize = decodeLz4(compressedBlobReader, decompressBlobArray, compressedBlobSize, uncompressedBlobSize, decompressBlobArray.decompressOffset)
						uncompressedBlobSize, err := lz4.UncompressBlockWithDict(src, decompressArray2, dict)

						log.Println(src)
						if err != nil {
							return nil, fmt.Errorf("failed to decode lz4 <2> in parseBinaryKv3Element: <%w>", err)
						}
						context.decompressOffset += uint32(uncompressedBlobSize)
						suboffset += uint32(uncompressedBlobSize)
						break
					}
				}
				return decompressArray, nil
			} else {
				panic("code me") /*
					if (uncompressedBlobReader) {//blobs have already been uncompressed
						let uncompressedBlobSize = uncompressedBlobSizeReader.getUint32();
						return uncompressedBlobReader.getBytes(uncompressedBlobSize);
					} else {
						//should not happend
						if (TESTING) {
							throw 'Missing reader';
						}
					}*/
			}
		}
	case kv3.DATA_TYPE_ARRAY:
		var count uint32
		err := binary.Read(quadReader, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to read count in parseBinaryKv3Element: <%w>", err)
		}

		elements := make([]kv3.Kv3Value, count)
		/*for (let i = 0; i < count; i++) {
			elements.push(readBinaryKv3Element(byteReader, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, decompressBlobArray, compressedBlobReader, uncompressedBlobReader, typeArray, valueArray, undefined, true, compressionFrameSize));
		}*/

		for i := uint32(0); i < count; i++ {
			value, err := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, compressedBlobReader, uncompressedBlobReader, typeReader, valueArray, elementType, false, compressionFrameSize)

			if err != nil {
				return nil, fmt.Errorf("call to parsebinarykv3element return an error, type %d, iteration %d: <%w>", elementType, i, err)
			}

			elements[i] = value
		}
		return elements, nil
	case kv3.DATA_TYPE_OBJECT:
		var count, nameId uint32
		err := binary.Read(quadReader, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to read count in parseBinaryKv3Element: <%w>", err)
		}

		element := kv3.NewKv3Element()
		for i := uint32(0); i < count; i++ {
			//let nameId = quadReader.getUint32();

			err := binary.Read(quadReader, binary.LittleEndian, &nameId)
			if err != nil {
				return nil, fmt.Errorf("failed to read nameId in parseBinaryKv3Element: <%w>", err)
			}

			elementType, err := readElementType(typeReader)
			if err != nil {
				return nil, fmt.Errorf("failed to read elementType in parseBinaryKv3Element: <%w>", err)
			}

			value, err := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, compressedBlobReader, uncompressedBlobReader, typeReader, valueArray, elementType, false, compressionFrameSize)
			if err != nil {
				return nil, fmt.Errorf("call to parsebinarykv3element return an error, type %d, iteration %d: <%w>", elementType, i, err)
			}

			element.AddAttribute(context.stringDictionary[nameId], value)
		}
		return element, nil
		/*
			count = quadReader.getUint32();
			//elements = new Kv3Element();
			elements = new Map();
			for (let i = 0; i < count; i++) {
				let nameId = quadReader.getUint32();
				let element = readBinaryKv3Element(byteReader, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, decompressBlobArray, compressedBlobReader, uncompressedBlobReader, typeArray, valueArray, undefined, false, compressionFrameSize);
				elements.set(nameId, element);
				//elements.setProperty(nameId, element);
			}
			return elements;*/
	case kv3.DATA_TYPE_TYPED_ARRAY:
		var count uint32
		err := binary.Read(quadReader, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to read count in parseBinaryKv3Element: <%w>", err)
		}

		subType, err := readElementType(typeReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read elementType, type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}

		elements := make([]kv3.Kv3Value, count)
		for i := uint32(0); i < count; i++ {
			//elements.push(readBinaryKv3Element(byteReader, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, decompressBlobArray, compressedBlobReader, uncompressedBlobReader, typeArray, valueArray, subType, true, compressionFrameSize));
			value, err := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, compressedBlobReader, uncompressedBlobReader, typeReader, valueArray, subType, true, compressionFrameSize)

			if err != nil {
				return nil, fmt.Errorf("call to parsebinarykv3element return an error, type %d, iteration %d: <%w>", elementType, i, err)
			}

			elements[i] = value
		}
		return elements, nil
	case kv3.DATA_TYPE_INT32:
		var value int32
		err := binary.Read(quadReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		return value, nil
	case kv3.DATA_TYPE_UINT32:
		var value uint32
		err := binary.Read(quadReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		return value, nil
	case kv3.DATA_TYPE_TRUE:
		return true, nil
	case kv3.DATA_TYPE_FALSE:
		return false, nil
	case kv3.DATA_TYPE_INT_ZERO:
		return 0, nil
	case kv3.DATA_TYPE_INT_ONE:
		return 1, nil
	case kv3.DATA_TYPE_DOUBLE_ZERO:
		return float64(0.0), nil
	case kv3.DATA_TYPE_DOUBLE_ONE:
		return float64(1.0), nil
	case kv3.DATA_TYPE_FLOAT:
		var value float32
		err := binary.Read(quadReader, binary.LittleEndian, &value)
		if err != nil {
			return nil, fmt.Errorf("failed to read value of type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}
		if isArray {
			return value, nil
		} else {
			valueArray = append(valueArray, value)
			return nil, nil
		}
	case kv3.DATA_TYPE_TYPED_ARRAY2:
		var count uint8
		err := binary.Read(context.reader, binary.LittleEndian, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to read count in parseBinaryKv3Element: <%w>", err)
		}

		log.Println("Count for DATA_TYPE_TYPED_ARRAY is", count)

		subType, err := readElementType(typeReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read elementType, type %d in parseBinaryKv3Element: <%w>", elementType, err)
		}

		elements := make([]kv3.Kv3Value, count)
		for i := uint8(0); i < count; i++ {
			//elements.push(readBinaryKv3Element(byteReader, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, decompressBlobArray, compressedBlobReader, uncompressedBlobReader, typeArray, valueArray, subType, true, compressionFrameSize));
			value, err := parseBinaryKv3Element(context, quadReader, eightReader, uncompressedBlobSizeReader, compressedBlobSizeReader, blobCount, decompressBlobBuffer, compressedBlobReader, uncompressedBlobReader, typeReader, valueArray, subType, true, compressionFrameSize)

			if err != nil {
				return nil, fmt.Errorf("call to parsebinarykv3element return an error, type %d, iteration %d: <%w>", elementType, i, err)
			}

			elements[i] = value
		}
		return elements, nil
	default:
		return nil, fmt.Errorf("unknown elementType %d in parsebinarykv3element", elementType)
	}
}

func readElementType(reader io.ReadSeeker) (byte, error) {
	var elementType byte
	err := binary.Read(reader, binary.LittleEndian, &elementType)
	if err != nil {
		return 0, fmt.Errorf("failed to read elementType in readElementType: <%w>", err)
	}

	if elementType == kv3.DATA_TYPE_RESOURCE {
		_, err := reader.Seek(1, io.SeekCurrent)
		if err != nil {
			return 0, fmt.Errorf("failed to seek in readElementType: <%w>", err)
		}
	}
	return elementType, nil
}

/*
	DATA_TYPE_NULL         = 0x01
	DATA_TYPE_BOOL         = 0x02
	DATA_TYPE_INT64        = 0x03
	DATA_TYPE_UINT64       = 0x04
	DATA_TYPE_DOUBLE       = 0x05
	DATA_TYPE_STRING       = 0x06
	DATA_TYPE_BLOB         = 0x07
	DATA_TYPE_ARRAY        = 0x08
	DATA_TYPE_OBJECT       = 0x09
	DATA_TYPE_TYPED_ARRAY  = 0x0A
	DATA_TYPE_INT32        = 0x0B
	DATA_TYPE_UINT32       = 0x0C
	DATA_TYPE_TRUE         = 0x0D
	DATA_TYPE_FALSE        = 0x0E
	DATA_TYPE_INT_ZERO     = 0x0F
	DATA_TYPE_INT_ONE      = 0x10
	DATA_TYPE_DOUBLE_ZERO  = 0x11
	DATA_TYPE_DOUBLE_ONE   = 0x12
	DATA_TYPE_FLOAT        = 0x13
	DATA_TYPE_BYTE         = 0x17
	DATA_TYPE_TYPED_ARRAY2 = 0x18
	DATA_TYPE_RESOURCE     = 0x86
*/

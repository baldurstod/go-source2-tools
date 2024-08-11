package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/repository"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
)

type parseContext struct {
	reader io.ReadSeeker
	file   *source2.File
	b      []byte
}

func newParseContext(repo string, filename string, b []byte) *parseContext {
	return &parseContext{
		reader: bytes.NewReader(b),
		file:   source2.NewFile(repo, filename),
		b:      b,
	}
}

func Parse(repo string, path string) (*source2.File, error) {
	reader := repository.GetRepository(repo)
	if reader == nil {
		return nil, errors.New("unknown repository")
	}

	b, err := reader.ReadFile(path)
	if err != nil {
		return nil, errors.New("unable to read file " + path)
	}

	context := newParseContext(repo, path, b)

	err = parseHeader(context)
	if err != nil {
		return nil, err
	}

	err = parseBlocks(context)
	if err != nil {
		return nil, err
	}

	return context.file, nil
}

func parseHeader(context *parseContext) error {
	//var test uint32
	var headerOffset uint32
	var resCount uint32
	var err error
	reader := context.reader

	err = binary.Read(reader, binary.LittleEndian, &context.file.FileLength)
	if err != nil {
		return fmt.Errorf("failed to read FileLength in parseHeader: <%w>", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &context.file.VersionMaj)
	if err != nil {
		return fmt.Errorf("failed to read VersionMaj in parseHeader: <%w>", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &context.file.VersionMin)
	if err != nil {
		return fmt.Errorf("failed to read VersionMin in parseHeader: <%w>", err)
	}

	currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
	err = binary.Read(reader, binary.LittleEndian, &headerOffset)
	if err != nil {
		return fmt.Errorf("failed to read headerOffset in parseHeader: <%w>", err)
	}
	headerOffset += uint32(currentPos)
	err = binary.Read(reader, binary.LittleEndian, &resCount)
	if err != nil {
		return fmt.Errorf("failed to read resCount in parseHeader: <%w>", err)
	}

	reader.Seek(int64(headerOffset), io.SeekStart)

	resType := make([]byte, 4)
	var resOffset uint32
	var resLength uint32
	for i := uint32(0); i < resCount; i++ {
		_, err := reader.Read(resType)
		if err != nil {
			return fmt.Errorf("failed to read block type in parseHeader: <%w>", err)
		}
		currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
		err = binary.Read(reader, binary.LittleEndian, &resOffset)
		if err != nil {
			return fmt.Errorf("failed to read resOffset in parseHeader: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &resLength)
		if err != nil {
			return fmt.Errorf("failed to read resLength in parseHeader: <%w>", err)
		}

		context.file.AddBlock(context.file.FileType, string(resType), resOffset+uint32(currentPos), resLength)
	}
	return nil
}

func parseBlocks(context *parseContext) error {
	for _, block := range context.file.BlocksArray {
		err := parseBlock(context, block)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseBlock(context *parseContext, block *source2.FileBlock) error {
	var err error
	switch block.ResType {
	case "RERL":
		err = parseRERL(context, block)
	case "DATA":
		err = parseBlock2(context, block)
	case "ANIM", "CTRL", "MRPH", "MDAT", "ASEQ", "AGRP", "PHYS", "LaCo":
		err = parseDATA(context, block)
	case "MBUF":
		err = parseVbib(context, block)
	default:
		log.Println("Unknown block type", block.ResType)
	}

	if err != nil {
		return fmt.Errorf("an error occured in parseblock: <%w>", err)
	}

	return nil
}

func parseRERL(context *parseContext, block *source2.FileBlock) error {
	var resOffset uint32
	var resCount uint32
	var strOffset int32
	var err error
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)

	err = binary.Read(reader, binary.LittleEndian, &resOffset)
	if err != nil {
		return fmt.Errorf("failed to read resOffset in parseRERL: <%w>", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &resCount)
	if err != nil {
		return fmt.Errorf("failed to read resCount in parseRERL: <%w>", err)
	}

	fileBlockRERL := source2.NewFileBlockRERL()
	block.Content = fileBlockRERL

	handle := make([]byte, 8)
	for i := uint32(0); i < resCount; i++ {
		context.reader.Seek(int64(block.Offset+resOffset+16*i), io.SeekStart)

		_, err := reader.Read(handle)
		if err != nil {
			return fmt.Errorf("failed to read handle in parseRERL: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &strOffset)
		if err != nil {
			return fmt.Errorf("failed to read strOffset in parseRERL: <%w>", err)
		}
		context.reader.Seek(int64(strOffset-4), io.SeekCurrent)
		filename, err := readNullString(reader)
		if err != nil {
			return fmt.Errorf("failed to read filename in parseRERL: <%w>", err)
		}

		fileBlockRERL.AddExternalFile(string(handle[:]), filename)

		//readHandle(reader)
	}
	/*
		function loadRerl(reader, block) {
			reader.seek(block.offset);
			var resOffset = reader.getInt32();// Seems to be always 0x00000008
			var resCount = reader.getInt32();
			block.externalFiles = {};
			block.externalFiles2 = [];

			reader.seek(block.offset + resOffset);


			for (var resIndex = 0; resIndex < resCount; resIndex++) {
				reader.seek(block.offset + resOffset + 16 * resIndex);
				var handle = readHandle(reader);//reader.getUint64(fieldOffset);
				var strOffset = reader.getInt32();
				reader.skip(strOffset - 4);
				var s = reader.getNullString();
				block.externalFiles[handle] = s;
				block.externalFiles2[resIndex] = s;
			}
		}
	*/

	return nil
}

func parseBlock2(context *parseContext, block *source2.FileBlock) error {
	if block.ResType != "DATA" {
		return errors.New("unsupported block type")
	}
	if block.Length <= 4 {
		return errors.New("can't determine data type")
	}

	var magic uint32
	context.reader.Seek(int64(block.Offset), io.SeekStart)
	switch magic {
	case 0x03564B56: // VKV3
		return parseDATAVKV3(context, block)
	case 0x4B563301: // kv31
		return parseDataKV3(context, block, 1)
	case 0x4B563302: // kv32 ?? new since wind ranger arcana
		return parseDataKV3(context, block, 2)
	case 0x4B563303: // KV3 v3 new since muerta
		return parseDataKV3(context, block, 3)
	case 0x4B563304: // KV3 v4 new since dota 7.33
		return parseDataKV3(context, block, 4)
	default:
		return parseData(context, block)
	}
}

func parseData(context *parseContext, block *source2.FileBlock) error {
	switch context.file.FileType {
	case "vcdlist":
		return parseVcdList(context, block)
	default:
		return errors.New("unsupported file type " + context.file.FileType)
	}
}

func parseDATA(context *parseContext, block *source2.FileBlock) error {
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)
	fileBlockDATA := source2.NewFileBlockDATA()
	block.Content = fileBlockDATA
	var magic uint32

	err := binary.Read(reader, binary.LittleEndian, &magic)
	if err != nil {
		return fmt.Errorf("failed to read magic in parseDATA: <%w>", err)
	}

	switch magic {
	case 0x03564B56: // VKV3
		return parseDATAVKV3(context, block)
	case 0x4B563301: // kv31
		return parseDataKV3(context, block, 1)
	case 0x4B563302: // kv32 ?? new since wind ranger arcana
		return parseDataKV3(context, block, 2)
	case 0x4B563303: // KV3 v3 new since muerta
		return parseDataKV3(context, block, 3)
	case 0x4B563304: // KV3 v4 new since dota 7.33
		return parseDataKV3(context, block, 4)
	default:
		log.Println("Unknown magic in parseDATA:", magic)
		/*	if (TESTING) {
			console.warn('Unknown block data type:', bytes);
		}*/
	}

	log.Println(fileBlockDATA, magic)
	return nil
	/*
		async loadData(reader, reference, block, introspection, parseVtex) {
			var bytes = reader.getUint32(block.offset);
			switch (bytes) {
				case 0x03564B56: // VKV3
					return loadDataVkv(reader, block);
				case 0x4B563301: // kv31
					return await loadDataKv3(reader, block, 1);
				case 0x4B563302: // kv32 ?? new since wind ranger arcana
					return await loadDataKv3(reader, block, 2);
				case 0x4B563303: // KV3 v3 new since muerta
					return await loadDataKv3(reader, block, 3);
				case 0x4B563304: // KV3 v4 new since dota 7.33
					return await loadDataKv3(reader, block, 4);
				default:
					if (TESTING) {
						console.warn('Unknown block data type:', bytes);
					}
			}
			if (!introspection || !introspection.structsArray) {
				if (parseVtex) {//TODO
					return loadDataVtex(reader, block);
				}
				return null;
			}
			block.structs = {};

			let structList = introspection.structsArray;
			var startOffset = block.offset;
			for (var structIndex = 0; structIndex < 1; structIndex++) {
				var struct = structList[structIndex];//introspection.firstStruct;
				block.structs[struct.name] = loadStruct(reader, reference, struct, block, startOffset, introspection, 0);
				startOffset += struct.discSize;
			}
			if (VERBOSE) {
				console.log(block.structs);
			}
		}
	*/
}

func parseDATAVKV3(context *parseContext, block *source2.FileBlock) error {
	log.Println("Code parseDATAVKV3", context, block)
	return nil
}

func parseDataKV3(context *parseContext, block *source2.FileBlock, version int) error {
	reader := context.reader
	reader.Seek(int64(block.Offset+4+16), io.SeekStart)

	var compressionMethod, singleByteCount, quadByteCount, eightByteCount uint32
	var dictionaryTypeLength, decodeLength, compressedLength uint32
	var blobCount, totalUncompressedBlobSize uint32
	var unknown5, unknown6 uint32
	var compressionDictionaryId, compressionFrameSize, unknown3, unknown4 uint16
	var err error
	compressedLength = block.Length

	err = binary.Read(reader, binary.LittleEndian, &compressionMethod)
	if err != nil {
		return fmt.Errorf("failed to read compressionMethod in parseDataKV3: <%w>", err)
	}
	if version >= 2 {
		err = binary.Read(reader, binary.LittleEndian, &compressionDictionaryId)
		if err != nil {
			return fmt.Errorf("failed to read compressionDictionaryId in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &compressionFrameSize)
		if err != nil {
			return fmt.Errorf("failed to read compressionFrameSize in parseDataKV3: <%w>", err)
		}
		//unknown1 = reader.getUint32();//0 or 0x40000000 depending on compression method
	}

	err = binary.Read(reader, binary.LittleEndian, &singleByteCount)
	if err != nil {
		return fmt.Errorf("failed to read singleByteCount in parseDataKV3: <%w>", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &quadByteCount)
	if err != nil {
		return fmt.Errorf("failed to read quadByteCount in parseDataKV3: <%w>", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &eightByteCount)
	if err != nil {
		return fmt.Errorf("failed to read eightByteCount in parseDataKV3: <%w>", err)
	}

	if version >= 2 {
		err = binary.Read(reader, binary.LittleEndian, &dictionaryTypeLength)
		if err != nil {
			return fmt.Errorf("failed to read dictionaryTypeLength in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &unknown3)
		if err != nil {
			return fmt.Errorf("failed to read unknown3 in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &unknown4)
		if err != nil {
			return fmt.Errorf("failed to read unknown4 in parseDataKV3: <%w>", err)
		}
	}

	err = binary.Read(reader, binary.LittleEndian, &decodeLength)
	if err != nil {
		return fmt.Errorf("failed to read decodeLength in parseDataKV3: <%w>", err)
	}

	if version >= 2 {
		err = binary.Read(reader, binary.LittleEndian, &compressedLength)
		if err != nil {
			return fmt.Errorf("failed to read compressedLength in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &blobCount)
		if err != nil {
			return fmt.Errorf("failed to read blobCount in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &totalUncompressedBlobSize)
		if err != nil {
			return fmt.Errorf("failed to read totalUncompressedBlobSize in parseDataKV3: <%w>", err)
		}
	}

	if version >= 4 {
		err = binary.Read(reader, binary.LittleEndian, &unknown5)
		if err != nil {
			return fmt.Errorf("failed to read unknown5 in parseDataKV3: <%w>", err)
		}
		err = binary.Read(reader, binary.LittleEndian, &unknown6)
		if err != nil {
			return fmt.Errorf("failed to read unknown6 in parseDataKV3: <%w>", err)
		}
	}

	var dst []byte

	var compressedBlobReader, uncompressedBlobReader io.ReadSeeker

	switch compressionMethod {
	case 0:
		dst = make([]byte, decodeLength)
		_, err := reader.Read(dst)
		if err != nil {
			return fmt.Errorf("failed to read datas in parseDataKV3 for compression method %d: <%w>", compressionMethod, err)
		}
	case 1:
		if blobCount > 0 {
			currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
			compressedBlobReader = bytes.NewReader(context.b[uint32(currentPos)+compressedLength:])
			//compressedBlobReader = new BinaryReader(reader, reader.tell() + compressedLength);
		}
		src := make([]byte, compressedLength)
		dst = make([]byte, decodeLength)
		_, err := reader.Read(src)
		if err != nil {
			return fmt.Errorf("failed to read datas in parseDataKV3 for compression method %d: <%w>", compressionMethod, err)
		}

		size, err := lz4.UncompressBlock(src, dst)
		if err != nil || uint32(size) != decodeLength {
			return fmt.Errorf("failed to decode lz4 in parseDataKV3 for compression method %d: <%w>", compressionMethod, err)
		}
	case 2: //new since spectre arcana
		decoder, _ := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
		currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
		src := context.b[uint32(currentPos) : uint32(currentPos)+compressedLength]
		sa, err := decoder.DecodeAll(src, nil)
		dst = sa[0:decodeLength]
		if err != nil {
			return fmt.Errorf("failed to decode zstd in parseDataKV3 for compression method %d: <%w>", compressionMethod, err)
		}
		if blobCount > 0 {
			uncompressedBlobReader = bytes.NewReader(sa[decodeLength:])
		}

	default:
		return fmt.Errorf("unknow compression method in parsedatakv3: %d", compressionMethod)
	}

	kv, err := ParseKv3(dst, version, singleByteCount, quadByteCount, eightByteCount, dictionaryTypeLength, blobCount, totalUncompressedBlobSize, compressedBlobReader, uncompressedBlobReader, compressionFrameSize)
	if err != nil {
		return fmt.Errorf("failed to parse kv3 in parseDataKV3: <%w>", err)
	}
	block.Content.(*source2.FileBlockDATA).KeyValue = kv

	//log.Println(compressionMethod, compressionDictionaryId, compressionFrameSize, singleByteCount, quadByteCount, eightByteCount, compressedLength)
	//log.Println(dictionaryTypeLength, unknown3, unknown4, decodeLength)
	//log.Println(compressedLength, blobCount, totalUncompressedBlobSize)

	/*
		async function loadDataKv3(reader, block, version) {
			const KV3_ENCODING_BLOCK_COMPRESSED = '\x46\x1A\x79\x95\xBC\x95\x6C\x4F\xA7\x0B\x05\xBC\xA1\xB7\xDF\xD2';
			const KV3_ENCODING_BLOCK_COMPRESSED_LZ4 = '\x8A\x34\x47\x68\xA1\x63\x5C\x4F\xA1\x97\x53\x80\x6F\xD9\xB1\x19';
			const KV3_ENCODING_BLOCK_COMPRESSED_UNKNOWN = '\x7C\x16\x12\x74\xE9\x06\x98\x46\xAF\xF2\xE6\x3E\xB5\x90\x37\xE7';

			let sa;
			let compressedBlobReader;
			let uncompressedBlobReader;
			switch (compressionMethod) {
				case 0:
					if (TESTING && version >= 2 && (compressionDictionaryId != 0 || compressionFrameSize != 0)) {
						throw 'Error compression method doesn\'t match';
					}
					sa = reader.getBytes(decodeLength);
					break;
				case 1:
					let buf = new ArrayBuffer(decodeLength);
					sa = new Uint8Array(buf);
					if (blobCount > 0) {
						compressedBlobReader = new BinaryReader(reader, reader.tell() + compressedLength);
					}
					decodeLz4(reader, sa, compressedLength, decodeLength);
					{
						if (blobCount > 0) {
							//SaveFile(new File([new Blob([sa])], 'decodeLz4_' + block.offset + '_' + block.length));
						}
						if (TESTING && (block.type == 'ANIM')) {
							//SaveFile(new File([new Blob([sa])], 'decodeLz4_block_ANIM_' + block.length + '_' + block.offset));
						}
					}
					break;
				case 2://new since spectre arcana
					//SaveFile(new File([new Blob([reader.getBytes(block.length, block.offset)])], 'block_' + block.offset + '_' + block.length));
					let compressedBytes = reader.getBytes(compressedLength);
					//SaveFile(new File([new Blob([compressedBytes])], 'block_' + block.offset + '_' + block.length));
					let decompressedBytes = await Zstd.decompress(compressedBytes);
					sa = new Uint8Array(new Uint8Array(decompressedBytes.buffer, 0, decodeLength));
					if (blobCount > 0) {
						uncompressedBlobReader = new BinaryReader(decompressedBytes, decodeLength);
					}
					//console.error(sa);
					//SaveFile(new File([new Blob([sa])], 'zstd'));

					break;
				default:
					throw 'Unknow kv3 compressionMethod ' + compressionMethod;
					break;
			}
			block.keyValue = BinaryKv3Loader.getBinaryKv3(version, sa, singleByteCount, quadByteCount, eightByteCount, dictionaryTypeLength, blobCount, totalUncompressedBlobSize, compressedBlobReader, uncompressedBlobReader, compressionFrameSize);
		}
	*/

	return nil
}

func readNullString(r io.ReadSeeker) (string, error) {
	// Probably not the fastest way to do that
	s := ""
	c := make([]byte, 1)

	for {
		_, err := r.Read(c)
		if err != nil {
			return s, fmt.Errorf("an error occured while reading null string: <%w>", err)
		}
		if c[0] == 0 {
			break
		} else {
			s += string(c[0])
		}
	}

	return s, nil
}

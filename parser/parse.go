package parser

import (
	"encoding/binary"
	"fmt"
	"github.com/baldurstod/go-source2-tools"
	"io"
	"log"
)

type parseContext struct {
	reader io.ReadSeeker
	file   *source2.File
}

func newParseContext(reader io.ReadSeeker) *parseContext {
	return &parseContext{
		reader: reader,
		file:   source2.NewFile(),
	}
}

func Parse(r io.ReadSeeker) (*source2.File, error) {
	context := newParseContext(r)

	log.Println("Start parsing file")
	err := parseHeader(context)
	if err != nil {
		return nil, err
	}

	err = parseBlocks(context)
	if err != nil {
		return nil, err
	}

	log.Println("End parsing file")
	return context.file, nil
}

func parseHeader(context *parseContext) error {
	//var test uint32
	var headerOffset uint32
	var resCount uint32
	reader := context.reader
	binary.Read(reader, binary.LittleEndian, &context.file.FileLength)
	binary.Read(reader, binary.LittleEndian, &context.file.VersionMaj)
	binary.Read(reader, binary.LittleEndian, &context.file.VersionMin)

	currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
	binary.Read(reader, binary.LittleEndian, &headerOffset)
	headerOffset += uint32(currentPos)
	binary.Read(reader, binary.LittleEndian, &resCount)

	reader.Seek(int64(headerOffset), io.SeekStart)

	resType := make([]byte, 4)
	var resOffset uint32
	var resLength uint32
	for i := uint32(0); i < resCount; i++ {
		_, err := reader.Read(resType)
		if err != nil {
			return fmt.Errorf("Failed to read block type in parseHeader: <%w>", err)
		}
		currentPos, _ := context.reader.Seek(0, io.SeekCurrent)
		binary.Read(reader, binary.LittleEndian, &resOffset)
		binary.Read(reader, binary.LittleEndian, &resLength)

		context.file.AddBlock(string(resType), resOffset+uint32(currentPos), resLength)
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

	log.Println(context, block)

	var err error
	switch block.ResType {
	case "RERL":
		err = parseRERL(context, block)
	case "DATA", "ANIM", "CTRL", "MRPH", "MDAT", "ASEQ", "AGRP", "PHYS", "LaCo":
		err = parseDATA(context, block)
	}

	if err != nil {
		return fmt.Errorf("An error occured in parseBlock: <%w>", err)
	}

	return nil
}

func parseRERL(context *parseContext, block *source2.FileBlock) error {
	var resOffset uint32
	var resCount uint32
	var strOffset int32
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)

	binary.Read(reader, binary.LittleEndian, &resOffset)
	binary.Read(reader, binary.LittleEndian, &resCount)

	log.Println(block.Offset, resOffset, resCount)

	fileBlockRERL := source2.NewFileBlockRERL(block)

	handle := make([]byte, 8)
	for i := uint32(0); i < resCount; i++ {
		context.reader.Seek(int64(block.Offset+resOffset+16*i), io.SeekStart)
		log.Println(block.Offset + resOffset + 16*i)

		_, err := reader.Read(handle)
		if err != nil {
			return fmt.Errorf("Failed to read handle in parseRERL: <%w>", err)
		}
		binary.Read(reader, binary.LittleEndian, &strOffset)
		context.reader.Seek(int64(strOffset-4), io.SeekCurrent)
		filename, err := readNullString(reader)
		if err != nil {
			return fmt.Errorf("Failed to read filename in parseRERL: <%w>", err)
		}

		fileBlockRERL.AddExternalFile(string(handle[:]), filename)

		//readHandle(reader)
	}
	log.Println(block.Offset, resOffset, resCount, fileBlockRERL)
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

func parseDATA(context *parseContext, block *source2.FileBlock) error {
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)
	fileBlockDATA := source2.NewFileBlockDATA(block)
	log.Println(context, block)
	var magic uint32

	binary.Read(reader, binary.LittleEndian, &magic)

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
	return nil
}
func parseDataKV3(context *parseContext, block *source2.FileBlock, version int) error {
	reader := context.reader
	reader.Seek(int64(block.Offset + 4 + 16), io.SeekStart)

	var compressionMethod, singleByteCount, quadByteCount, eightByteCount uint32
	var compressionDictionaryId, compressionFrameSize uint16
	compressedLength := block.Length

	binary.Read(reader, binary.LittleEndian, &compressionMethod)
	if (version >= 2) {
		binary.Read(reader, binary.LittleEndian, &compressionDictionaryId)
		binary.Read(reader, binary.LittleEndian, &compressionFrameSize)
		//unknown1 = reader.getUint32();//0 or 0x40000000 depending on compression method
	}

	binary.Read(reader, binary.LittleEndian, &singleByteCount)
	binary.Read(reader, binary.LittleEndian, &quadByteCount)
	binary.Read(reader, binary.LittleEndian, &eightByteCount)

	log.Println(compressionMethod, compressionDictionaryId, compressionFrameSize, singleByteCount, quadByteCount, eightByteCount, compressedLength)

/*
	format 	 := make([]byte, 16)
	_, err := reader.Read(format)
	if err != nil {
		return fmt.Errorf("Failed to read format in parseDATAVKV3: <%w>", err)
	}
	log.Println(string(format[:]))*/

	/*
	var compressionMethod uint32
	var resCount uint32
	var strOffset int32
	reader := context.reader
	reader.Seek(int64(block.Offset), io.SeekStart)

	binary.Read(reader, binary.LittleEndian, &resOffset)*/

	/*
		async function loadDataKv3(reader, block, version) {
			const KV3_ENCODING_BLOCK_COMPRESSED = '\x46\x1A\x79\x95\xBC\x95\x6C\x4F\xA7\x0B\x05\xBC\xA1\xB7\xDF\xD2';
			const KV3_ENCODING_BLOCK_COMPRESSED_LZ4 = '\x8A\x34\x47\x68\xA1\x63\x5C\x4F\xA1\x97\x53\x80\x6F\xD9\xB1\x19';
			const KV3_ENCODING_BLOCK_COMPRESSED_UNKNOWN = '\x7C\x16\x12\x74\xE9\x06\x98\x46\xAF\xF2\xE6\x3E\xB5\x90\x37\xE7';

			reader.seek(block.offset);

			let method = 1;

			reader.skip(4);
			let format = reader.getString(16);
			let compressionMethod = reader.getUint32();
			let compressionDictionaryId;
			let compressionFrameSize;
			let dictionaryTypeLength, unknown3, unknown4, blobCount = 0, totalUncompressedBlobSize;
			let unknown5, unknown6;
			if (version >= 2) {
				compressionDictionaryId = reader.getUint16();
				compressionFrameSize = reader.getUint16();
				//unknown1 = reader.getUint32();//0 or 0x40000000 depending on compression method
			}
			let singleByteCount = reader.getUint32();//skip this many bytes ?????
			let quadByteCount = reader.getUint32();
			let eightByteCount = reader.getUint32();
			let compressedLength = block.length;
			if (version >= 2) {
				dictionaryTypeLength = reader.getUint32();
				unknown3 = reader.getUint16();
				unknown4 = reader.getUint16();
				if (false && TESTING) {
					console.log(dictionaryTypeLength, unknown3, unknown4, block);
				}
			}

			var decodeLength = reader.getUint32();
			if (version >= 2) {
				compressedLength = reader.getUint32();
				blobCount = reader.getUint32();
				totalUncompressedBlobSize = reader.getUint32();
			}

			if (version >= 4) {
				unknown5 = reader.getUint32();
				unknown6 = reader.getUint32();
			}

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
			return s, fmt.Errorf("An error occured while reading null string: <%w>", err)
		}
		if c[0] == 0 {
			break
		} else {
			s += string(c[0])
		}
	}

	log.Println(s)
	return s, nil
}

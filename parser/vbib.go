package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/baldurstod/go-source2-tools"
	"github.com/baldurstod/go-source2-tools/vertex"
)

const (
	VERTEX_HEADER_SIZE    = 24
	INDEX_HEADER_SIZE     = 24
	DESC_HEADER_SIZE      = 56
	DESC_HEADER_NAME_SIZE = 36
)

func parseVbib(context *parseContext, block *source2.FileBlock) error {
	return nil
	blockReader := bytes.NewReader(context.b[block.Offset:])

	fileBlockMBUF := source2.NewFileBlockMBUF()
	block.Data = fileBlockMBUF

	var vertexOffset, vertexCount, indexOffset, indexCount int32
	err := binary.Read(blockReader, binary.LittleEndian, &vertexOffset)
	if err != nil {
		return fmt.Errorf("failed to read vertex offset: <%w>", err)
	}

	err = binary.Read(blockReader, binary.LittleEndian, &vertexCount)
	if err != nil {
		return fmt.Errorf("failed to read vertex count: <%w>", err)
	}

	err = binary.Read(blockReader, binary.LittleEndian, &indexOffset)
	if err != nil {
		return fmt.Errorf("failed to read vertex offset: <%w>", err)
	}
	indexOffset += 8

	err = binary.Read(blockReader, binary.LittleEndian, &indexCount)
	if err != nil {
		return fmt.Errorf("failed to read vertex count: <%w>", err)
	}

	log.Println(vertexOffset, vertexCount, indexOffset, indexCount)
	/*
	   block.vertices = [];
	   block.indices = [];
	*/
	var bytesPerVertex, headerOffset, headerCount, dataOffset, dataLength int32
	for i := int32(0); i < vertexCount; i++ { // header size: 24 bytes
		_, err := blockReader.Seek(int64(vertexOffset+i*VERTEX_HEADER_SIZE), io.SeekStart)
		if err != nil {
			return fmt.Errorf("failed to seek vertex %d: <%w>", i, err)
		}

		s1 := vertex.NewVertex()

		err = binary.Read(blockReader, binary.LittleEndian, &s1.VertexCount)
		if err != nil {
			return fmt.Errorf("failed to read vertex count: <%w>", err)
		}

		err = binary.Read(blockReader, binary.LittleEndian, &bytesPerVertex)
		if err != nil {
			return fmt.Errorf("failed to read bytes per vertex: <%w>", err)
		}

		err = binary.Read(blockReader, binary.LittleEndian, &headerOffset)
		if err != nil {
			return fmt.Errorf("failed to read header offset: <%w>", err)
		}

		err = binary.Read(blockReader, binary.LittleEndian, &headerCount)
		if err != nil {
			return fmt.Errorf("failed to read header count: <%w>", err)
		}

		err = binary.Read(blockReader, binary.LittleEndian, &dataOffset)
		if err != nil {
			return fmt.Errorf("failed to read data offset: <%w>", err)
		}

		err = binary.Read(blockReader, binary.LittleEndian, &dataLength)
		if err != nil {
			return fmt.Errorf("failed to read data lebgth: <%w>", err)
		}

		vertexDataSize := vertexCount * bytesPerVertex

		vertexReader := blockReader

		if vertexDataSize != dataLength {
			panic("TODO: MeshoptDecoder")
			/*
				let vertexBuffer = new Uint8Array(new ArrayBuffer(vertexDataSize));
				MeshoptDecoder.decodeVertexBuffer(vertexBuffer, s1.vertexCount, s1.bytesPerVertex, new Uint8Array(reader.buffer.slice(s1.dataOffset, s1.dataOffset + s1.dataLength)));
				//SaveFile('sa.obj', new Blob([vertexBuffer]));
				vertexReader = new BinaryReader(vertexBuffer);
				s1.dataOffset = 0;		*/
		}

		log.Println(s1.VertexCount, bytesPerVertex, headerOffset, headerCount, dataOffset, dataLength, vertexDataSize, vertexReader)
		panic("stop")
		/*

			let vertexDataSize = s1.vertexCount * s1.bytesPerVertex;

			let vertexReader = reader;
			if (vertexDataSize != s1.dataLength) {
				let vertexBuffer = new Uint8Array(new ArrayBuffer(vertexDataSize));
				MeshoptDecoder.decodeVertexBuffer(vertexBuffer, s1.vertexCount, s1.bytesPerVertex, new Uint8Array(reader.buffer.slice(s1.dataOffset, s1.dataOffset + s1.dataLength)));
				//SaveFile('sa.obj', new Blob([vertexBuffer]));
				vertexReader = new BinaryReader(vertexBuffer);
				s1.dataOffset = 0;
			}

			s1.headers = [];
			for (var j = 0; j < s1.headerCount; j++) { // header size: 24 bytes
				var header = {};
				var headerOffset = s1.headerOffset + j * DESC_HEADER_SIZE;
				reader.seek(headerOffset);
				header.name = reader.getNullString();
				reader.seek(headerOffset + DESC_HEADER_NAME_SIZE);
				header.type = reader.getUint32();
				header.offset = reader.getUint32();

				s1.headers.push(header);
			}

			s1.vertices = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_POSITION);
			s1.normals = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_NORMAL);
			s1.tangents = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_TANGENT);
			s1.coords = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_COORD);
			s1.boneIndices = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_BONE_INDICE);
			s1.boneWeight = new ArrayBuffer(s1.vertexCount * BYTES_PER_VERTEX_BONE_WEIGHT);

			let s1Vertices = new Float32Array(s1.vertices);
			let s1Normals = new Float32Array(s1.normals);
			let s1Tangents = new Float32Array(s1.tangents);
			let s1Coords = new Float32Array(s1.coords);
			let s1BoneIndices = new Float32Array(s1.boneIndices);
			let s1BoneWeight = new Float32Array(s1.boneWeight);
			for (var vertexIndex = 0; vertexIndex < s1.vertexCount; vertexIndex++) {
				vertexReader.seek(s1.dataOffset + vertexIndex * s1.bytesPerVertex);
				var vertex = {};

				var positionFilled = false;//TODOv3: remove this
				var normalFilled = false;
				var tangentFilled = false;
				var texCoordFilled = false;
				var blendIndicesFilled = false;
				var blendWeightFilled = false;
				for (var headerIndex = 0; headerIndex < s1.headers.length; headerIndex++) {
					var headerName = s1.headers[headerIndex].name;
					var headerType = s1.headers[headerIndex].type;
					let tempValue;// = vec4.create();//TODO: optimize


					vertexReader.seek(s1.dataOffset + vertexIndex * s1.bytesPerVertex + s1.headers[headerIndex].offset);
					switch (headerType) {
						case DXGI_FORMAT_R32G32B32A32_FLOAT:
							tempValue = vec4.create();//TODO: optimize
							tempValue[0] = vertexReader.getFloat32();
							tempValue[1] = vertexReader.getFloat32();
							tempValue[2] = vertexReader.getFloat32();
							tempValue[3] = vertexReader.getFloat32();
							break;
						case DXGI_FORMAT_R32G32B32_FLOAT:// 3 * float32
							tempValue = vec3.create();//TODO: optimize
							tempValue[0] = vertexReader.getFloat32();
							tempValue[1] = vertexReader.getFloat32();
							tempValue[2] = vertexReader.getFloat32();
							break;
						case DXGI_FORMAT_R16G16B16A16_SINT:
							tempValue = vec4.create();//TODO: optimize
							tempValue[0] = vertexReader.getInt16();
							tempValue[1] = vertexReader.getInt16();
							tempValue[2] = vertexReader.getInt16();
							tempValue[3] = vertexReader.getInt16();
							break;
						case DXGI_FORMAT_R32G32_FLOAT:// 2 * float32
							tempValue = vec2.create();//TODO: optimize
							tempValue[0] = vertexReader.getFloat32();
							tempValue[1] = vertexReader.getFloat32();
							break;
						case DXGI_FORMAT_R8G8B8A8_UNORM:
							tempValue = vec4.create();//TODO: optimize
							tempValue[0] = vertexReader.getUint8() / 255;
							tempValue[1] = vertexReader.getUint8() / 255;
							tempValue[2] = vertexReader.getUint8() / 255;
							tempValue[3] = vertexReader.getUint8() / 255;
							//vertexReader.getUint8();
							break;
						case DXGI_FORMAT_R8G8B8A8_UINT:// 4 * uint8
							tempValue = vec4.create();//TODO: optimize
							tempValue[0] = vertexReader.getUint8();
							tempValue[1] = vertexReader.getUint8();
							tempValue[2] = vertexReader.getUint8();
							tempValue[3] = vertexReader.getUint8();
							break;
						case DXGI_FORMAT_R16G16_FLOAT:// 2 * float16
							tempValue = vec2.create();//TODO: optimize
							tempValue[0] = vertexReader.getFloat16();
							tempValue[1] = vertexReader.getFloat16();
							break;
						case DXGI_FORMAT_R16G16_SNORM://New with battlepass 2022
							tempValue = vec2.create();//TODO: optimize
							tempValue[0] = sNormUint16(vertexReader.getInt16());
							tempValue[1] = sNormUint16(vertexReader.getInt16());
							break;
						case DXGI_FORMAT_R16G16_SINT:
							tempValue = vec2.create();//TODO: optimize
							tempValue[0] = vertexReader.getInt16();
							tempValue[1] = vertexReader.getInt16();
							break;
						case DXGI_FORMAT_R32_FLOAT	:// single float32 ??? new in half-life Alyx
							tempValue = [];
							tempValue[0] = vertexReader.getFloat32();
							break;
						case DXGI_FORMAT_R32_UINT: // single uint32 ??? new since DOTA2 2023_08_30
							tempValue = [];
							tempValue[0] = vertexReader.getUint32();
							//console.log(tempValue[0]);
							break;
						default:
							//TODO add types when needed. see DxgiFormat.js
							console.error('Warning: unknown type ' + headerType + ' for value ' + headerName);
							tempValue = vec4.create();//TODO: optimize
							tempValue[0] = 0;
							tempValue[1] = 0;
							tempValue[2] = 0;
							tempValue[3] = 0;
					}

					switch (headerName) {
						case 'POSITION':
							s1Vertices.set(tempValue, vertexIndex * VERTEX_POSITION_LEN);
							positionFilled = true;
							break;
						case 'NORMAL':
							s1Normals.set(tempValue, vertexIndex * VERTEX_NORMAL_LEN);//TODOv3
							normalFilled = true;
							break;
						case 'TANGENT':
							s1Tangents.set(tempValue, vertexIndex * VERTEX_NORMAL_LEN);//TODOv3
							tangentFilled = true;
							break;
						case 'TEXCOORD':
							if (!texCoordFilled) {//TODO: handle 2 TEXCOORD
								let test = vec2.clone(tempValue);//todov3: fixme see //./Alyx/models/props_industrial/hideout_doorway.vmdl_c
								s1Coords.set(test/*tempValue* /, vertexIndex * VERTEX_COORD_LEN);
								texCoordFilled = true;
							}
							break;
						case 'BLENDINDICES':
							/*s1.boneIndices.push(tempValue[0]);
							s1.boneIndices.push(tempValue[1]);
							s1.boneIndices.push(tempValue[2]);
							s1.boneIndices.push(tempValue[3]);* /
							s1BoneIndices.set(tempValue, vertexIndex * VERTEX_BONE_INDICE_LEN);
							blendIndicesFilled = true;
							break;
						case 'BLENDWEIGHT':
							/*s1.boneWeight.push(tempValue[0]);
							s1.boneWeight.push(tempValue[1]);
							s1.boneWeight.push(tempValue[2]);
							s1.boneWeight.push(tempValue[3]);* /
							//vec4.scale(tempValue, tempValue, 1 / 255.0);
							s1BoneWeight.set(tempValue, vertexIndex * VERTEX_BONE_WEIGHT_LEN);
							blendWeightFilled = true;
							break;
						//TODOv3: add "texcoord" lowercase maybe a z- tex coord ?
					}
				}

				//TODOv3: remove this
				if (!positionFilled) {
					/*s1.vertices.push(0);
					s1.vertices.push(0);
					s1.vertices.push(0);* /
					s1Vertices.set(defaultValuesPosition, vertexIndex * VERTEX_POSITION_LEN);
				}
				if (!normalFilled) {
					/*s1.normals.push(0);
					s1.normals.push(0);
					s1.normals.push(0);* /
					s1Normals.set(defaultValuesNormal, vertexIndex * VERTEX_NORMAL_LEN);
				}
				if (!tangentFilled) {
					s1Tangents.set(defaultValuesTangent, vertexIndex * VERTEX_TANGENT_LEN);
				}
				if (!texCoordFilled) {
					/*s1.coords.push(0);
					s1.coords.push(0);* /
					s1Coords.set(defaultValuesCoord, vertexIndex * VERTEX_COORD_LEN);
				}
				if (!blendIndicesFilled) {
					/*s1.boneIndices.push(0);
					s1.boneIndices.push(0);
					s1.boneIndices.push(0);
					s1.boneIndices.push(0);* /
					s1BoneIndices.set(defaultValuesBoneIndice, vertexIndex * VERTEX_BONE_INDICE_LEN);
				}
				if (!blendWeightFilled) {
					/*s1.boneWeight.push(255);
					s1.boneWeight.push(0);
					s1.boneWeight.push(0);
					s1.boneWeight.push(0);* /
					s1BoneWeight.set(defaultValuesBoneWeight, vertexIndex * VERTEX_BONE_WEIGHT_LEN);
				}

			}
			if (VERBOSE) {
				console.log(s1.normals[0], s1.normals[1], s1.normals[2]);
			}
			block.vertices.push(s1);
		*/
	}
	/*
	   //console.log(block.vertices);

	   for (var i = 0; i < indexCount; i++) { // header size: 24 bytes

	   		reader.seek(indexOffset + i * INDEX_HEADER_SIZE);
	   		var s2 = {};
	   		s2.indexCount = reader.getInt32();
	   		s2.bytesPerIndex = reader.getInt32();
	   		s2.headerOffset = reader.tell() + reader.getInt32();
	   		s2.headerCount = reader.getInt32();
	   		s2.dataOffset = reader.tell() + reader.getInt32();
	   		s2.dataLength = reader.getInt32();

	   		let indexDataSize = s2.indexCount * s2.bytesPerIndex;

	   		let indexReader = reader;
	   		if (indexDataSize != s2.dataLength) {
	   			let indexBuffer = new Uint8Array(new ArrayBuffer(indexDataSize));
	   			MeshoptDecoder.decodeIndexBuffer(indexBuffer, s2.indexCount, s2.bytesPerIndex, new Uint8Array(reader.buffer.slice(s2.dataOffset, s2.dataOffset + s2.dataLength)));
	   			indexReader = new BinaryReader(indexBuffer);
	   			s2.dataOffset = 0;
	   		}

	   		s2.indices = new ArrayBuffer(s2.indexCount * BYTES_PER_INDEX);
	   		let s2Indices = new Uint32Array(s2.indices);
	   		for (var indicesIndex = 0; indicesIndex < s2.indexCount; indicesIndex++) {
	   			indexReader.seek(s2.dataOffset + indicesIndex * s2.bytesPerIndex);
	   			var vertex = {};
	   			//s2.indices.push(indexReader.getUint16());
	   			if (s2.bytesPerIndex == 2) {
	   				s2Indices[indicesIndex] = indexReader.getUint16();
	   			} else {
	   				s2Indices[indicesIndex] = indexReader.getUint32();
	   			}
	   		}

	   		block.indices.push(s2);
	   	}
	*/
	panic("stop")
	return nil
}

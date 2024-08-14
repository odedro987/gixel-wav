package wav

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

const (
	RiffId = "RIFF"
	WaveId = "WAVE"
	FmtId  = "fmt "
	DataId = "data"

	AudiometadataPCM = 1
)

type ChunkLookupTable = map[string]int

type WavMetadata struct {
	AudioFormat   int
	ChannelNum    int
	SampleRate    int
	ByteRate      int
	BlockAlign    int
	BitsPerSample int
}

type Signal struct {
	SampleRate int
	Samples    []float32
}

func Decode(path string) (*WavMetadata, *Signal, error) {
	var metadata WavMetadata

	lookupTable, err := analyzeChunks(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed analyzing chunks: %w", err)
	}
	fmt.Println(lookupTable)

	err = decodeFMTChunk(path, lookupTable[FmtId], &metadata)
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding FMT chunk: %w", err)
	}

	var signal Signal
	err = decodeData(path, lookupTable[DataId], metadata, &signal)
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding data: %w", err)
	}

	return &metadata, &signal, nil
}

func analyzeChunks(path string) (ChunkLookupTable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed opening file: %w", err)
	}
	defer file.Close()

	err = decodeRIFFChunk(file)
	if err != nil {
		return nil, fmt.Errorf("failed decoding RIFF chunk: %w", err)
	}

	lookupTable := make(ChunkLookupTable)
	lookupTable[RiffId] = 0

	var offset = 12 // RIFF id, size and WAVE id

	for {
		_, err = file.Seek(int64(offset), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed seeking offset %d: %w", offset, err)
		}

		chunkIdBytes := make([]byte, 4)
		_, err = file.Read(chunkIdBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed reading bytes: %w", err)
		}

		chunkSizeBytes := make([]byte, 4)
		_, err = file.Read(chunkSizeBytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed reading bytes: %w", err)
		}

		var chunkSize = int(binary.LittleEndian.Uint32(chunkSizeBytes)) + 8 // 8 bytes are id + size
		lookupTable[string(chunkIdBytes)] = offset

		offset += chunkSize
	}

	return lookupTable, nil
}

func decodeData(path string, offset int, metadata WavMetadata, signal *Signal) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed seeking chunk: %s, offset: %d, err: %w", DataId, offset, err)
	}

	chunkIdBytes := make([]byte, 4)
	_, err = file.Read(chunkIdBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	if string(chunkIdBytes) != DataId {
		return fmt.Errorf("file is not WAV metadata: chunk id(%s) is not %s", string(chunkIdBytes), DataId)
	}

	chunkSizeBytes := make([]byte, 4)
	_, err = file.Read(chunkSizeBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	var chunkSize = int(binary.LittleEndian.Uint32(chunkSizeBytes))
	var bytesPerSample = metadata.BitsPerSample / 8

	var sampleNum = chunkSize / (bytesPerSample * metadata.ChannelNum)
	var normFactor = 1.0 / (math.Pow(2, float64(metadata.BitsPerSample-1)) - 1)
	fmt.Println("number of samples", sampleNum)
	samples := make([]float32, sampleNum)

	var bytes = make([]byte, chunkSize)
	_, err = file.Read(bytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}

	// TODO: support all BitsPerSample
	for i := 0; i < sampleNum; i++ {
		var offset = i * bytesPerSample * metadata.ChannelNum
		test := int(int16(binary.LittleEndian.Uint16(bytes[offset : offset+bytesPerSample])))
		samples[i] = float32(test) * float32(normFactor)
	}

	signal.SampleRate = metadata.SampleRate
	signal.Samples = samples

	return nil
}

func decodeFMTChunk(path string, offset int, metadata *WavMetadata) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed seeking chunk: %s, offset: %d, err: %w", FmtId, offset, err)
	}

	chunkIdBytes := make([]byte, 4)
	_, err = file.Read(chunkIdBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	if string(chunkIdBytes) != FmtId {
		return fmt.Errorf("file is not WAV metadata: chunk id is not %s", FmtId)
	}

	chunkSizeBytes := make([]byte, 4)
	_, err = file.Read(chunkSizeBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}

	audioFormatBytes := make([]byte, 2)
	_, err = file.Read(audioFormatBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.AudioFormat = int(binary.LittleEndian.Uint16(audioFormatBytes))

	channelNumBytes := make([]byte, 2)
	_, err = file.Read(channelNumBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.ChannelNum = int(binary.LittleEndian.Uint16(channelNumBytes))

	sampleRateBytes := make([]byte, 4)
	_, err = file.Read(sampleRateBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.SampleRate = int(binary.LittleEndian.Uint32(sampleRateBytes))

	byteRateBytes := make([]byte, 4)
	_, err = file.Read(byteRateBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.ByteRate = int(binary.LittleEndian.Uint32(byteRateBytes))

	blockAlignBytes := make([]byte, 2)
	_, err = file.Read(blockAlignBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.BlockAlign = int(binary.LittleEndian.Uint16(blockAlignBytes))

	bitPerSampleBytes := make([]byte, 2)
	_, err = file.Read(bitPerSampleBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	metadata.BitsPerSample = int(binary.LittleEndian.Uint16(bitPerSampleBytes))

	return nil
}

func decodeRIFFChunk(file *os.File) error {
	chunkIdBytes := make([]byte, 4)
	_, err := file.Read(chunkIdBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	if string(chunkIdBytes) != RiffId {
		return fmt.Errorf("file is not WAV metadata: chunk id is not \"RIFF\"")
	}

	chunkSizeBytes := make([]byte, 4)
	_, err = file.Read(chunkSizeBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}

	metadataBytes := make([]byte, 4)
	_, err = file.Read(metadataBytes)
	if err != nil {
		return fmt.Errorf("failed reading bytes: %w", err)
	}
	if string(metadataBytes) != WaveId {
		return fmt.Errorf("file is not WAV metadata: metadata is not \"WAVE\"")
	}

	return nil
}

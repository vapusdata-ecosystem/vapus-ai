package appcl

import (
	"encoding/json"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	filetools "github.com/vapusdata-ecosystem/vapusai/core/tools/files"
	"google.golang.org/protobuf/types/known/structpb"
)

type VapusStreamOpts struct {
	Event       mpb.VapusStreamEvents
	Content     any
	Reason      mpb.EOSReasons
	ResError    error
	ContentType mpb.ContentFormats
	ContentId   string
}

func BuildNabhikTaskStreamChunk(opts *VapusStreamOpts) *mpb.NabhikTaskStreamChunk {
	var contents = ""
	var convertToBytes = false
	data := &mpb.VapusContentObject{
		ContentType: opts.ContentType,
		Final:       &mpb.VapusEOL{},
	}
	switch opts.ContentType {
	case mpb.ContentFormats_JSON:
		mappertype, ok := opts.Content.(*mpb.Mapper)
		if ok {
			opts.Content = map[string]any{
				mappertype.Key: mappertype.Value,
			}
		}
		contents, err := filetools.GenericMarshaler(opts.Content, opts.ContentType.String())
		if err != nil {
			convertToBytes = true
		} else {
			data.Content = string(contents)
		}
	case mpb.ContentFormats_YAML:
		contents, err := filetools.GenericMarshaler(opts.Content, opts.ContentType.String())
		if err != nil {
			convertToBytes = true
		} else {
			data.Content = string(contents)
		}
	case mpb.ContentFormats_MAP:
		vv, ok := opts.Content.(*structpb.Struct)
		if !ok {
			convertToBytes = true
		} else {
			contents, err := filetools.GenericMarshaler(vv, mpb.ContentFormats_JSON.String())
			if err != nil {
				convertToBytes = true
			} else {
				data.Content = string(contents)
			}
		}
	case mpb.ContentFormats_PLAIN_TEXT:
		vv, ok := opts.Content.(string)
		if !ok {
			convertToBytes = true
		} else {
			data.Content = vv
		}
	case mpb.ContentFormats_FORMATTED_CONTENT:
		vv, ok := opts.Content.(string)
		if !ok {
			convertToBytes = true
		} else {
			data.Content = vv
		}
	case mpb.ContentFormats_CLICKSET:
		vv, ok := opts.Content.(string)
		if !ok {
			convertToBytes = true
		} else {
			data.Content = vv
		}
	default:
		contents, err := filetools.GenericMarshaler(opts.Content, opts.ContentType.String())
		if err != nil {
			convertToBytes = true
		} else {
			data.Content = string(contents)
		}
	}
	if convertToBytes {
		vv, ok := opts.Content.([]byte)
		if !ok {
			bbytes, err := json.Marshal(opts.Content)
			if err != nil {
				data.Content = "error while processing content"
			} else {
				data.Content = string(bbytes)
			}
		} else {
			data.Content = string(vv)
		}
	}
	if opts.Event == mpb.VapusStreamEvents_DATA ||
		opts.Event == mpb.VapusStreamEvents_DATASET_START ||
		opts.Event == mpb.VapusStreamEvents_START ||
		opts.Event == mpb.VapusStreamEvents_DATASET_END ||
		opts.Event == mpb.VapusStreamEvents_REASONINGS ||
		opts.Event == mpb.VapusStreamEvents_DATASET ||
		opts.Event == mpb.VapusStreamEvents_DATETIME ||
		opts.Event == mpb.VapusStreamEvents_FILE_DATA ||
		opts.Event == mpb.VapusStreamEvents_STATE ||
		opts.Event == mpb.VapusStreamEvents_CHAT_OVERFLOWEN ||
		opts.Event == mpb.VapusStreamEvents_TASK_CREATED ||
		opts.Event == mpb.VapusStreamEvents_RESPONSE_ID {
		data.Final = nil
	} else {
		data.Final.Metadata = func() string {
			if opts.ResError == nil {
				return string(contents)
			}
			return opts.ResError.Error()
		}()
		data.Final.Reason = opts.Reason
	}
	if opts.ContentId != "" {
		data.ContentId = opts.ContentId
	}
	obj := &mpb.NabhikTaskStreamChunk{
		EventAt: dmutils.GetEpochTime(),
		Event:   opts.Event,
		Data:    data,
	}
	if opts.Event == mpb.VapusStreamEvents_FILE_DATA {
		fd, ok := opts.Content.(*mpb.FileData)
		if !ok {
			opts.Event = mpb.VapusStreamEvents_DATA
		} else {
			obj.Data.Content = ""
			obj.Files = fd
		}
	}
	return obj
}

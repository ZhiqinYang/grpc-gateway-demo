package runtime

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// 数据替换
func init() {
	runtime.HTTPError = HTTPErrorHandler
}

type errorBody struct {
	Msg  string          `protobuf:"bytes,1,name=msg" json:"msg"`
	Code int32           `protobuf:"varint,2,name=code" json:"code"`
	Data json.RawMessage `protobuf:"bytes,3,rep,name=data" json:"data,omitempty"`
}

func HTTPErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	body := &errorBody{
		Msg:  s.Message(),
		Code: int32(s.Code()),
		Data: []byte(s.Message()),
	}

	buf, merr := json.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message: %v", merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, mux, md)
	handleForwardResponseTrailerHeader(w, md)
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
	handleForwardResponseTrailer(w, md)
}

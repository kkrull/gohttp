package playground

import "io"

type StatelessOptionController struct {
}

func (controller *StatelessOptionController) Options(client io.Writer, target string) {
	//TODO KDK: Copy-paste the response from capability/, then see if it can be extracted to an http.Response implementation
	panic("implement me")
}

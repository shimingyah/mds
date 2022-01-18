package server

import "net/http"

// Status return mds current status
func (m *MDS) Status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

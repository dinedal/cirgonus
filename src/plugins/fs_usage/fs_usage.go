package fs_usage

/*
// int statfs(const char *path, struct statfs *buf);

#include <sys/vfs.h>
#include <stdlib.h>
#include <assert.h>

struct statfs* go_statfs(const char *path) {
  struct statfs *fsinfo;
  fsinfo = malloc(sizeof(struct statfs));
  assert(fsinfo != NULL);
  statfs(path, fsinfo);
  return fsinfo;
}
*/
import "C"

import (
	"fmt"
	"logger"
	"unsafe"
)

func GetMetric(params interface{}, log *logger.Logger) interface{} {
	path := C.CString(params.(string))
	stat := C.go_statfs(path)

	log.Log("debug", fmt.Sprintf("blocks size on %s: %v", string(*path), stat.f_bsize))
	log.Log("debug", fmt.Sprintf("blocks total on %s: %v", string(*path), stat.f_blocks))
	log.Log("debug", fmt.Sprintf("blocks free on %s: %v", string(*path), stat.f_bfree))

	defer C.free(unsafe.Pointer(stat))
	defer C.free(unsafe.Pointer(path))

	return [3]uint64{
		(uint64(((float64(stat.f_blocks - stat.f_bfree)) / float64(stat.f_blocks)) * 100)),
		(uint64(stat.f_bfree) * uint64(stat.f_bsize)),
		(uint64(stat.f_blocks) * uint64(stat.f_bsize)),
	}
}

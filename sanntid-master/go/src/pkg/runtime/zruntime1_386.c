// auto generated by go tool dist
// goos=linux goarch=386

#include "runtime.h"
void
runtime·GOMAXPROCS(int32 n, int32 ret)
{
#line 940 "/tmp/bindist204690179/go/src/pkg/runtime/runtime1.goc"

	ret = runtime·gomaxprocsfunc(n);
	FLUSH(&ret);
}
void
runtime·NumCPU(int32 ret)
{
#line 944 "/tmp/bindist204690179/go/src/pkg/runtime/runtime1.goc"

	ret = runtime·ncpu;
	FLUSH(&ret);
}
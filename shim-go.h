#ifndef __YAJL_SHIM_GO_H__
#define __YAJL_SHIM_GO_H__

YAJL_API char * yajl_val_get_string(yajl_val val);
YAJL_API char * yajl_val_get_number(yajl_val val);
YAJL_API size_t yajl_val_get_array_len(yajl_val val);
YAJL_API size_t yajl_val_get_object_len(yajl_val val);
YAJL_API const char *yajl_val_get_object_key(yajl_val val, size_t index);
YAJL_API yajl_val yajl_val_get_object_value(yajl_val val, size_t index);
YAJL_API yajl_val yajl_val_get_array_value(yajl_val val, size_t index);

#endif

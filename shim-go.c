#include "api/yajl_tree.h"

//since golang has no union, and access union from golang is quite difficult
//this shim is use to support helper function to access the yajl_val

char *
yajl_val_get_string(yajl_val val)
{
	return val->u.string;
}

char *
yajl_val_get_number(yajl_val val)
{
	return val->u.number.r;
}

size_t
yajl_val_get_array_len(yajl_val val)
{
	return val->u.array.len;
}

size_t
yajl_val_get_object_len(yajl_val val)
{
	return val->u.object.len;
}

const char *
yajl_val_get_object_key(yajl_val val, size_t index)
{
	return val->u.object.keys[index];
}

yajl_val
yajl_val_get_object_value(yajl_val val, size_t index)
{
	return val->u.object.values[index];
}

yajl_val
yajl_val_get_array_value(yajl_val val, size_t index)
{
	return val->u.array.values[index];
}
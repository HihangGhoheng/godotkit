package gdk_converter

func StringToPointer(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func PointerToString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func IntToPointer(i int) *int {
	return &i
}

func PointerToInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func Int8ToPointer(i int8) *int8 {
	return &i
}

func PointerToInt8(i *int8) int8 {
	if i == nil {
		return int8(0)
	}
	return *i
}

func Int16ToPointer(i int16) *int16 {
	return &i
}

func PointerToInt16(i *int16) int16 {
	if i == nil {
		return int16(0)
	}
	return *i
}

func Int32ToPointer(i int32) *int32 {
	return &i
}

func PointerToInt32(i *int32) int32 {
	if i == nil {
		return int32(0)
	}
	return *i
}

func Int64ToPointer(i int64) *int64 {
	return &i
}

func PointerToInt64(i *int64) int64 {
	if i == nil {
		return int64(0)
	}
	return *i
}

func UintToPointer(u uint) *uint {
	return &u
}

func PointerToUint(u *uint) uint {
	if u == nil {
		return uint(0)
	}
	return *u
}

func Uint8ToPointer(u uint8) *uint8 {
	return &u
}

func PointerToUint8(u *uint8) uint8 {
	if u == nil {
		return uint8(0)
	}
	return *u
}

func Uint64ToPointer(u uint64) *uint64 {
	return &u
}

func PointerToUint64(u *uint64) uint64 {
	if u == nil {
		return uint64(0)
	}
	return *u
}

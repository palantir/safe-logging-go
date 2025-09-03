package innera

type InnerTestStructWithSafeField struct {
	SafeField *string `safelogging:"@Safe"`
}

type InnerTestStructWithUnsafeField struct {
	UnsafeField *string `safelogging:"@Unsafe"`
}

type InnerTestStructWithDoNotLogField struct {
	DoNotLogField *string `safelogging:"@DoNotLog"`
}

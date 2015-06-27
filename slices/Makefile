types:= int int64 int32 uint uint64 uint32 rune byte bool

all: $(types)
	@echo string

$(types): string.go
	@echo $@
	@cp $< $@.go; sed -i "s/String/$(shell echo $@ | sed 's/\b[a-z]/\U&/')/g" $@.go; sed -i "s/string/$@/g" $@.go

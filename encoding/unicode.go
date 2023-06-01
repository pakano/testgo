package encoding

/*
Unicode：所有字符都是两个字节，对于英文字符，高字节为0，低字节与ASCII码相同
Unicode的实现方式不同于编码方式。一个字符的Unicode编码确定。
但是在实际传输过程中，由于不同系统平台的设计不一定一致，以及出于节省空间的目的，对Unicode编码的实现方式有所不同。
Unicode的实现方式称为Unicode转换格式（Unicode Transformation Format，简称为UTF）。
*/

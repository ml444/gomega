package subscribe

//func (s *SubConfig) UnMarshalRequest(data []byte) (proto.Message, error) {
//	inT := reflect.TypeOf(s.Request).Elem()
//	in := reflect.New(inT).Interface().(proto.Message)
//	err := proto.Unmarshal(data, in)
//	if err != nil {
//		return nil, err
//	}
//	return in, nil
//}
//
//func (s SubConfig) NewResponse() proto.Message {
//	T := reflect.TypeOf(s.Response).Elem()
//	return reflect.New(T).Interface().(proto.Message)
//}

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domain

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4a0f95aaDecode20242FIGHTCLUBDomain(in *jlexer.Lexer, out *UserResponce) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "rating":
			out.Rating = float64(in.Float64())
		case "avatar":
			out.Avatar = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "sex":
			out.Sex = string(in.String())
		case "birthDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthdate).UnmarshalJSON(data))
			}
		case "guestCount":
			out.GuestCount = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain(out *jwriter.Writer, in UserResponce) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"rating\":"
		out.RawString(prefix[1:])
		out.Float64(float64(in.Rating))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"sex\":"
		out.RawString(prefix)
		out.String(string(in.Sex))
	}
	{
		const prefix string = ",\"birthDate\":"
		out.RawString(prefix)
		out.Raw((in.Birthdate).MarshalJSON())
	}
	{
		const prefix string = ",\"guestCount\":"
		out.RawString(prefix)
		out.Int(int(in.GuestCount))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserResponce) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserResponce) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserResponce) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserResponce) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain1(in *jlexer.Lexer, out *UserDataResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "uuid":
			out.Uuid = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "score":
			out.Score = float64(in.Float64())
		case "avatar":
			out.Avatar = string(in.String())
		case "sex":
			out.Sex = string(in.String())
		case "guestCount":
			out.GuestCount = int(in.Int())
		case "birthdate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthdate).UnmarshalJSON(data))
			}
		case "isHost":
			out.IsHost = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain1(out *jwriter.Writer, in UserDataResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"uuid\":"
		out.RawString(prefix[1:])
		out.String(string(in.Uuid))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"score\":"
		out.RawString(prefix)
		out.Float64(float64(in.Score))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"sex\":"
		out.RawString(prefix)
		out.String(string(in.Sex))
	}
	{
		const prefix string = ",\"guestCount\":"
		out.RawString(prefix)
		out.Int(int(in.GuestCount))
	}
	{
		const prefix string = ",\"birthdate\":"
		out.RawString(prefix)
		out.Raw((in.Birthdate).MarshalJSON())
	}
	{
		const prefix string = ",\"isHost\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsHost))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserDataResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserDataResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserDataResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserDataResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain1(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain2(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.UUID = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "score":
			out.Score = float64(in.Float64())
		case "avatar":
			out.Avatar = string(in.String())
		case "sex":
			out.Sex = string(in.String())
		case "guestCount":
			out.GuestCount = int(in.Int())
		case "birthDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthdate).UnmarshalJSON(data))
			}
		case "isHost":
			out.IsHost = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain2(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.UUID))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"score\":"
		out.RawString(prefix)
		out.Float64(float64(in.Score))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"sex\":"
		out.RawString(prefix)
		out.String(string(in.Sex))
	}
	{
		const prefix string = ",\"guestCount\":"
		out.RawString(prefix)
		out.Int(int(in.GuestCount))
	}
	{
		const prefix string = ",\"birthDate\":"
		out.RawString(prefix)
		out.Raw((in.Birthdate).MarshalJSON())
	}
	{
		const prefix string = ",\"isHost\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsHost))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain2(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain3(in *jlexer.Lexer, out *UpdateUserRegion) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "regionName":
			out.RegionName = string(in.String())
		case "startVisitedDate":
			out.StartVisitedDate = string(in.String())
		case "endVisitedDate":
			out.EndVisitedDate = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain3(out *jwriter.Writer, in UpdateUserRegion) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"regionName\":"
		out.RawString(prefix[1:])
		out.String(string(in.RegionName))
	}
	{
		const prefix string = ",\"startVisitedDate\":"
		out.RawString(prefix)
		out.String(string(in.StartVisitedDate))
	}
	{
		const prefix string = ",\"endVisitedDate\":"
		out.RawString(prefix)
		out.String(string(in.EndVisitedDate))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UpdateUserRegion) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UpdateUserRegion) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UpdateUserRegion) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UpdateUserRegion) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain3(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain4(in *jlexer.Lexer, out *SessionData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain4(out *jwriter.Writer, in SessionData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SessionData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SessionData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SessionData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SessionData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain4(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain5(in *jlexer.Lexer, out *GetAllUsersResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]*UserDataResponse, 0, 8)
					} else {
						out.Users = []*UserDataResponse{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *UserDataResponse
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(UserDataResponse)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Users = append(out.Users, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain5(out *jwriter.Writer, in GetAllUsersResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"users\":"
		out.RawString(prefix[1:])
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Users {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetAllUsersResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetAllUsersResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetAllUsersResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetAllUsersResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain5(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain6(in *jlexer.Lexer, out *CSRFTokenResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "csrf_token":
			out.Token = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain6(out *jwriter.Writer, in CSRFTokenResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"csrf_token\":"
		out.RawString(prefix[1:])
		out.String(string(in.Token))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CSRFTokenResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CSRFTokenResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CSRFTokenResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CSRFTokenResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain6(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain7(in *jlexer.Lexer, out *AuthResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "session_id":
			out.SessionId = string(in.String())
		case "user":
			(out.User).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain7(out *jwriter.Writer, in AuthResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"session_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.SessionId))
	}
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix)
		(in.User).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AuthResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AuthResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AuthResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AuthResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain7(l, v)
}
func easyjson4a0f95aaDecode20242FIGHTCLUBDomain8(in *jlexer.Lexer, out *AuthData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "username":
			out.Username = string(in.String())
		case "email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson4a0f95aaEncode20242FIGHTCLUBDomain8(out *jwriter.Writer, in AuthData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AuthData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AuthData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncode20242FIGHTCLUBDomain8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AuthData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AuthData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecode20242FIGHTCLUBDomain8(l, v)
}

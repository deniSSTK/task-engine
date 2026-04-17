package authRepo

import "go.uber.org/zap/zapcore"

func (u *UpdateUser) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", u.Id.String())

	if u.Name != nil {
		enc.AddString("name", *u.Name)
	}

	if u.SecondName != nil {
		if *u.SecondName == nil {
			enc.AddString("second_name", "nil")
		} else {
			enc.AddString("second_name", **u.SecondName)
		}
	}

	return nil
}

package types

type Password string

func (p Password) String() string { return string(p) }

func (p Password) SecurityString() string { return "--------" }

func (p *Password) Decode(method string) (raw string, err error) { return }

func (p Password) Encode(method string) (security string, err error) { return }

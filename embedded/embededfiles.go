package embedded

import "embed"

//go:embed lua_injections/*
var LuaInjections embed.FS

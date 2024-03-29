package main

import "testing"

func assertParseLine(t *testing.T, src string, expect asmLine) {
	t.Logf("Calling parseLine(%q)...", src)
	actual := parseLine(src)
	if actual == nil {
		t.Error("Did not expected nil result")
		return
	}
	if actual.label != expect.label {
		t.Errorf("Labels don't match! Expected %q, got %q", expect.label, actual.label)
	}
	if actual.text != expect.text {
		t.Errorf("Texts don't match! Expected %q, got %q", expect.text, actual.text)
	}
	if actual.comment != expect.comment {
		t.Errorf("Comments don't match! Expected %q, got %q", expect.comment, actual.comment)
	}
}

func TestParseLine(t *testing.T) {
	assertParseLine(t, "", asmLine{"", "", ""})
	assertParseLine(t, ":", asmLine{"", "", ""})
	assertParseLine(t, ";", asmLine{"", "", ""})
	assertParseLine(t, "::;", asmLine{"", ":", ""})
	assertParseLine(t, "    label:    ", asmLine{"label", "", ""})
	assertParseLine(t, "   label:    ;    and a comment    ", asmLine{"label", "", "and a comment"})
	assertParseLine(t, "   mov rax, rbx    ;    and a comment    ", asmLine{"", "mov rax, rbx", "and a comment"})
	assertParseLine(t, "   mov rax, rbx    ;      ", asmLine{"", "mov rax, rbx", ""})
	assertParseLine(t, "   mov   rax,rbx    ;      ", asmLine{"", "mov rax, rbx", ""})
	assertParseLine(t, "   mov   rax,   rbx    ", asmLine{"", "mov rax, rbx", ""})
	assertParseLine(t, "   lbl  :mov   rax,   rbx    ", asmLine{"lbl", "mov rax, rbx", ""})
	assertParseLine(t, "   lbl  :mov   rax,   rbx    \n  \n", asmLine{"lbl", "mov rax, rbx", ""})
	assertParseLine(t, "mov eax, fs:[eax + 0x30]", asmLine{"", "mov eax, fs:[eax + 0x30]", ""})

	assertParseLine(t, "lbl db data", asmLine{"lbl", "db data", ""})
	assertParseLine(t, " dw   data", asmLine{"", "dw data", ""})
	assertParseLine(t, "lbl   dd   data;", asmLine{"lbl", "dd data", ""})
	assertParseLine(t, "lbl dq data", asmLine{"lbl", "dq data", ""})
	assertParseLine(t, "lbl ddq data", asmLine{"lbl", "ddq data", ""})
	assertParseLine(t, "lbl do data", asmLine{"lbl", "do data", ""})
	assertParseLine(t, "lbl dt data", asmLine{"lbl", "dt data", ""})
	assertParseLine(t, "add data", asmLine{"", "add data", ""})

	assertParseLine(t, "lbl: ;:comment", asmLine{"lbl", "", ":comment"})
	assertParseLine(t, "data ;:comment", asmLine{"", "data", ":comment"})
	assertParseLine(t, ";fake_lbl:comment", asmLine{"", "", "fake_lbl:comment"})

	assertParseLine(t, `db "not: a label"`, asmLine{"", `db "not: a label"`, ""})
	assertParseLine(t, `db ";not a comment" ;real comment `, asmLine{"", `db ";not a comment"`, "real comment"})
}

func assertNoquotes(t *testing.T, in, rep, expect string) {
	actual := noquotes(in, rep)
	if actual != expect {
		t.Errorf("noquotes(%q, %q): Expected %q, got %q", in, rep, expect, actual)
	}
}

func TestNoquotes(t *testing.T) {
	assertNoquotes(t, "", "x", "")
	assertNoquotes(t, "normal string", "x", "normal string")
	assertNoquotes(t, "something 'quoted'", "x", "something xxxxxxxx")
	assertNoquotes(t, `something "quoted"`, "x", "something xxxxxxxx")
	assertNoquotes(t, `something "quoted" and 'again'`, "x", "something xxxxxxxx and xxxxxxx")
	assertNoquotes(t, `more "qu'ot" marks`, "x", "more xxxxxxx marks")
	assertNoquotes(t, `more "qu'ot" 'ma"r"ks`, "x", "more xxxxxxx 'maxxxks")
	assertNoquotes(t, `more "qu'ot" 'marks`, "x", "more xxxxxxx 'marks")
	assertNoquotes(t, `more "qu'ot" 'marks`, " ", "more         'marks")
	assertNoquotes(t, `more "qu'ot" 'marks`, "", "more  'marks")

	assertNoquotes(t, "в кавычках 'юникод'", "x", "в кавычках xxxxxxxx")
}

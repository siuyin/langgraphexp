# My experiment is reproducing LangGraph capability in Go.

With "google.golang.org/genai" v1.12.0 I'm getting failures with chat.SendMessageStream

The failure is intermmitent and appears to be related to bug: https://github.com/googleapis/go-genai/issues/310.

A problematic prompt is "mocha please".
```
goroutine 1 [running]:
main.chatStream-range2(0xc000036770?, {0x0?, 0x0?})
        /home/siuyin/exp/langgraph/go_equiv/main.go:62 +0xfe
google.golang.org/genai.(*Chat).SendStream.func1-range1(0xc0000a95c0, {0x0?, 0x0?})
        /home/siuyin/go/pkg/mod/google.golang.org/genai@v1.12.0/chats.go:151 +0x1a5
google.golang.org/genai.Models.generateContentStream.iterateResponseStream[...].func5()
        /home/siuyin/go/pkg/mod/google.golang.org/genai@v1.12.0/api_client.go:258 +0x2c2
google.golang.org/genai.(*Chat).SendStream.func1(0xc000589b20)
        /home/siuyin/go/pkg/mod/google.golang.org/genai@v1.12.0/chats.go:140 +0x11c
main.chatStream(0xc000187680)
        /home/siuyin/exp/langgraph/go_equiv/main.go:58 +0x2e9
main.main()
        /home/siuyin/exp/langgraph/go_equiv/main.go:37 +0x159
exit status 2
```

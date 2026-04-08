# CLAUDE.md

이 파일은 Claude Code (claude.ai/code)가 이 저장소에서 작업할 때 참고하는 가이드입니다.
공통 가이드는 [AGENTS.md](AGENTS.md)를 참조하세요.

## Claude 전용 지침

- 한국어로 응답하세요.
- 코드 변경 시 반드시 `go test ./... -v`로 검증 후 완료 처리.
- TDD 우선: 새 기능 추가 시 테스트 먼저 작성 → 실패 확인 → 구현 → 통과 확인.
- 예제 바이너리 빌드 시 `go build -o /tmp/bin ./example/<name>/` 사용 (디렉토리 이름 충돌 방지).

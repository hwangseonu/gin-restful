# CLAUDE.md

이 파일은 Claude Code (claude.ai/code)가 이 저장소에서 작업할 때 참고하는 가이드입니다.

## 프로젝트 개요

`gin-restful`은 Gin 웹 프레임워크 위에 Flask-RESTful 스타일의 추상화를 제공하는 Go 라이브러리입니다. 구현한 인터페이스만 라우트로 등록되며, 제네릭 헬퍼로 타입 안전한 JSON 바인딩을 지원합니다.

## 명령어

```bash
go build ./...              # 라이브러리 빌드
go vet ./...                # 린트
go test ./... -v            # 테스트 실행
go run ./example/...        # 예제 서버 실행 (:8080)
go build -o /tmp/sample ./example/  # 예제 바이너리 빌드 (디렉토리 충돌 방지)
```

## 아키텍처

### 핵심 파일

- **`resource.go`** — 6개 독립 인터페이스 (`Lister`, `Getter`, `Poster`, `Putter`, `Patcher`, `Deleter`). 필요한 것만 구현하면 해당 HTTP 메서드만 라우트로 등록됨.
- **`api.go`** — `API` struct와 `AddResource`. 타입 어서션으로 resource가 구현한 인터페이스를 감지하고 Gin 라우트를 등록.
- **`bind.go`** — `Bind[T]` 제네릭 헬퍼. `c.ShouldBindJSON`을 감싸서 타입 안전한 요청 바디 파싱 제공.
- **`errors.go`** — `HTTPError` 타입과 `Abort()` 헬퍼. Flask-RESTful의 `abort()` 역할.
- **`handler.go`** — `makeHandler` (응답 디스패처), `handleError` (HTTPError 감지), `normalizePath` (경로 `//` 방지).

### 설계 패턴

- **인터페이스 분리**: 모놀리식 Resource 인터페이스 대신, HTTP 메서드별 독립 인터페이스. 부분 CRUD 자연 지원.
- **하이브리드 타입 안전성**: 메서드 시그니쳐는 `(any, int, error)` 반환으로 통일, 요청 바디는 `Bind[T]`로 컴파일 타임 타입 검증.
- **에러 처리**: `errors.As`로 `*HTTPError`를 감지하여 상태코드+JSON 응답 자동 생성. 일반 error는 500.
- **204 No Content**: status가 204이면 본문 없이 `c.Status()` 호출.

## 패키지

Go 패키지 이름은 `gin_restful`이며, 관례적으로 `restful`로 import합니다.

# Do Memory Viewer

React 기반 Do Memory 웹 뷰어입니다.

## 실행 방법

```bash
cd .do/viewer
npm install
npm run dev
```

브라우저에서 http://localhost:3777 으로 접속합니다.

## 사전 요구사항

Go Worker API가 localhost:3778에서 실행 중이어야 합니다:

```bash
cd .do/worker
go run .
```

## 구조

```
src/
├── api/
│   └── client.ts        # Worker API 클라이언트
├── components/
│   ├── Layout.tsx       # 공통 레이아웃
│   ├── Sidebar.tsx      # 사이드바 네비게이션
│   └── Timeline.tsx     # 타임라인 컴포넌트
├── pages/
│   ├── Dashboard.tsx    # 대시보드 (통계, 최근 활동)
│   ├── Sessions.tsx     # 세션 목록 및 상세
│   ├── Observations.tsx # 관찰 검색 및 필터링
│   ├── Plans.tsx        # 플랜 목록 및 내용
│   └── Reports.tsx      # 일별 요약 리포트
├── App.tsx
├── main.tsx
└── index.css
```

## 기술 스택

- React 19
- TypeScript
- Vite
- TailwindCSS
- React Router 7

## API 엔드포인트

Worker API (localhost:3778):

- `GET /api/sessions` - 세션 목록
- `GET /api/sessions/:id` - 세션 상세
- `GET /api/observations` - 관찰 목록
- `GET /api/observations/search?q=` - 관찰 검색
- `GET /api/plans` - 플랜 목록
- `GET /api/summaries?days=7` - 일별 요약
- `GET /health` - 헬스 체크

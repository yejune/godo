---
name: do-framework-electron
description: >
  Electron 33+ 데스크톱 앱 개발 전문가 - Main/Renderer 프로세스 아키텍처,
  IPC 통신, auto-update, Electron Forge 및 electron-builder로 패키징,
  보안 모범 사례涵盖. cross-platform 데스크톱 애플리케이션 구축,
  네이티브 OS 통합 구현, 또는 배포를 위한 Electron 앱 패키징 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.0.0"
  category: "framework"
  status: "active"
  updated: "2026-01-10"
  modularized: "false"
  tags: "electron, desktop, cross-platform, nodejs, chromium, ipc, auto-update, electron-builder, electron-forge"
  context7-libraries: "/electron/electron, /electron/forge, /electron-userland/electron-builder"
  related-skills: "do-lang-typescript, do-domain-frontend, do-lang-javascript"

---

# Electron 33+ 데스크톱 개발

## Quick Reference

Electron 33+ 데스크톱 앱 개발 전문가로 web 기술로 cross-platform 데스크톱 애플리케이션을 구축할 수 있습니다.

자동 트리거: electron.vite.config.ts 또는 electron-builder.yml 파일이 감지되는 Electron 프로젝트, 데스크톱 앱 개발 요청, IPC 통신 패턴 구현

핵심 기능:

Electron 33 플랫폼:

- Chromium 130 렌더링 엔진으로 모던 웹 기능
- Node.js 20.18 런타임으로 네이티브 시스템 액세스
- Main 프로세스에서 네이티브 ESM 지원
- GPU 가속 그래픽을 위한 WebGPU API 지원

프로세스 아키텍처:

- Main 프로세스는 애플리케이션당 단일 인스턴스로 전체 Node.js 액세스를 가짐
- Renderer 프로세스는 샌드박스된 환경에서 웹 콘텐츠를 표시
- Preload 스크립트는 Main과 Renderer를 제어된 API 노출로 브리지
- Utility 프로세스는 UI를 차단하지 않고 백그라운드 작업을 처리

IPC 통신:

- 요청-응답 통신을 위한 타입 안전 invoke/handle 패턴
- Main 프로세스 기능을 안전하게 Renderer에 노출하기 위한 contextBridge API
- Main에서 Renderer로 푸시 알림을 위한 이벤트 기반 메시징

Auto-Update 지원:

- GitHub 및 S3 게시와 electron-updater 통합
- 더 작은 다운로드 크기의 차등 업데이트
- 업데이트 알림 및 설치 관리

패키징 옵션:

- 통합 빌드 도구와 플그인 에코시템을 위한 Electron Forge
- 유연한 multi-platform 배포를 위한 electron-builder

보안 기능:

- 프로토타입 오염을 방지하기 위한 contextIsolation
- Renderer 프로세스 격리를 위한 Sandbox 강제
- Content Security Policy 구성
- IPC handler를 위한 입력 검증 패턴

### Project Initialization

새 Electron 애플리케이션을 생성하려면 create-electron-app 명령을 vite-typescript 템플릿과 함께 실행하세요. 패키징을 위해 electron-builder를 dev dependency로 설치하세요. auto-update 기능을 위해 electron-updater를 runtime dependency로 설치하세요.

상세한 명령과 구성은 reference.md Quick Commands 섹션을 참조하세요.

---

## Implementation Guide

### Project Structure

권장 디렉토리 레아웃:

소스 디렉토리는 4개의 주요 하위 디렉토리를 포함해야 합니다:

Main 디렉토리: main 프로세스 진입점, 도메인별로 구성된 IPC handler, 비즈니스 로직 서비스, 윈도 관리 모듈을 포함합니다.

Preload 디렉토리: preload 스크립트 진입점과 Main과 Renderer를 연결하는 API 정의를 포함합니다.

Renderer 디렉토리: React, Vue, 또는 Svelte로 구축된 웹 애플리케이션을 포함합니다. HTML 진입점과 Vite 구성을 포함합니다.

Shared 디렉토리: Main과 Renderer 프로세스 간에 공유되는 TypeScript 타입과 상수를 포함합니다.

프로젝트 루트에는 build 구성을 위한 electron.vite.config.ts, 패키징 옵션을 위한 electron-builder.yml, 앱 아이콘 및 자산을 위한 resources 디토리가 포함되어야 합니다.

### Main Process Setup

애플리케이션 수명주기 관리:

main 프로세스 초기화는 특정 순서를 따릅니다. 먼저 app.enableSandbox()를 호출하여 전역적으로 sandbox를 활성화하여 모든 renderer 프로세스가 격리된 환경에서 실행되도록 하세요. 앱이케이션의 단일 인스턴스 실행을 방지하기 위해 app.requestSingleInstanceLock()을 요청하세요.

창 생성은 app ready 이벤트가 발생한 후에 이루어져야 합니다. BrowserWindow를 webPreferences로 구성하세요. contextIsolation을 활성화하고, nodeIntegration을 비활성화하고, sandbox를 활성화하고, webSecurity를 활성화하세요. preload 스크립트 경로를 설정하여 Renderer에 안전한 API를 노출하세요.

macOS 동작을 위해: 창이 없을 때 dock 아이콘을 클릭하면 창을 재생성합니다. 다른 플랫폼에서는 모든 창이 닫히면 애플리케이션을 종료합니다.

구현 예제는 examples.md Main Process Entry Point 섹션을 참조하세요.

### 타입 안전 IPC 통신

IPC 타입 정의:

채널 이름을 페이로드 타입에 매핑하는 interface를 정의하세요. file operations, window operations, storage operations과 같이 도메인별로 그룹화하세요. 이를 통해 Main process handler와 Renderer invocation 모두에 대한 타입 검사가 가능합니다.

Main Process Handler 등록:

공유 타입을 임포트하는 전용 모듈에서 IPC handler를 등록하세요. 각 handler는 Zod와 같은 스키마 검증 라이브러리를 사용하여 처리 전에 입력을 검증해야 합니다. 요청-응답 패턴을 위해 ipcMain.handle을 사용하고 구조화된 결과를 반환하세요.

Preload 스크립트 구현:

각 채널의 ipcRenderer.invoke 호출을 래핑는 API 객체를 생성하세요. contextBridge.exposeInMainWorld를 사용하여 이 API를 window.electronAPI로 Renderer에서 사용 가능하게 만드세요. 이벤트 리스너러를 위한 cleanup 함수를 포함하여 메모리 누수를 방지하세요.

완전한 IPC 구현 패턴은 examples.md Type-Safe IPC Implementation 섹션을 참조하세요.

### 보안 모범 사례

필수 보안 설정:

모든 BrowserWindow는 4가지 필수 webPreferences를 구성해야 합니다. contextIsolation은 항상 활성화되어야 Electron 내부를 Renderer 코드로부터 보호해야 합니다. nodeIntegration은 Renderer 프로세스에서 항상 비활성화되어야 합니다. sandbox는 프로세스 수준 격리를 위해 항상 활성화되어야 합니다. webSecurity는 same-origin 정책 강제를 유지하기 위해 절대 비활성화되어야 합니다.

Content Security Policy:

webRequest.onHeadersReceived를 사용하여 session-level CSP 헤더를 구성하세요. default-src를 self로, script-src를 unsafe-inline 없이 self로, connect-src를 허용된 API 도메인으로 제한하세요. 이는 XSS 공격과 무단 리소스 로딩을 방지합니다.

입력 검증:

모든 IPC handler는 처리 전에 입력을 검증해야 합니다. 경로 순회 공격을 방지하기 위해 상위 디토리에 대한 참조를 거부하세요. 파일명에서 예약된 문자를 유효성 검사합니다. 파일 액세스를 구현할 때 허용된 디렉토리에 대한 allowlist를 사용하세요.

보안 구현 세부 정보는 reference.md Security Best Practices 섹션을 참조하세요.

### Auto-Update 구현

업데이트 서비스 아키텍처:

electron-updater 라이프사이클을 관리하는 UpdateService 클래스를 생성하세요. UI 알림을 활성화하기 위해 main window 참조로 초기화하세요. 사용자 대역 제어를 위해 autoDownload를 false로 설정하세요.

이벤트 처리:

update-available 이벤트를 처리하여 사용자에게 다운로드 확인을 요청하세요. download-progress 이벤트를 처리하여 진행률 표시기기를 표시하세요. update-downloaded 이벤트를 처리하여 재시작을 요청하세요.

사용자 알림 패턴:

시스템 다이얼로그를 사용하여 업데이트 가능 및 다운로드 완료 시 사용자에게 프롬프트를 제공하세요. Renderer에 이벤트를 보내어 앱 내 알림을 표시하세요. 즉시 및 지연 설치 모두 지원하세요.

완전한 업데이트 서비스 구현은 examples.md Auto-Update Integration 섹션을 참조하세요.

### 앱 패키징

Electron Builder 구성:

역-domain 표기법으로 appId를 설정하여 플랫폼 등록을 위합니다. 시스템 UI에 표시될 productName을 지정하세요. macOS, Windows, Linux용 platform별 target을 설정하세요.

macOS 구성:

App Store 분류를 위한 category를 설정하세요. notarization을 위해 hardenedRuntime을 활성화하고 entitlements를 구성하세요. x64와 arm64 아키텍처를 모두 타겟팅하는 universal build를 구성하세요.

Windows 구성:

실행 파일과 설치 프로그램을 위한 icon 경로를 지정하세요. 설치 디렉토리 선택을 위한 NSIS 설치 옵션을 구성하세요. 적절한 해시 알고리즘으로 code signing을 설정하세요.

Linux 구성:

데스크톱 환경 통합을 위한 category를 설정하세요. AppImage(universal 배포)와 deb/rpm(패키지 매니저 설치)를 포함한 다중 target을 설정하세요.

완전한 구성 예제는 reference.md Configuration 섹션을 참조하세요.

---

## Advanced Patterns

포괄적인 문서는 다음을 포함합니다:

Window 상태 지속성:

- 세션 간 창 위치 및 크기 저장 및 복원
- 여러 디스플레이 및 디스플레이 변경 처리
- 최대화 및 전체화 화면 상태 관리

Multi-Window 관리:

- 적절한 parent-child 관계로 이차 창 생성
- 여러 창 간 상태 공유
- 창 수명주기 이벤트 조율

시스템 트레이 및 네이티브 메뉴:

- context menus가 있는 시스템 트레이 아이콘 생성
- 키보드 단축키가 있는 애플리케이션 메뉴 구성
- macOS 및 Windows용 플랫폼별 메뉴 패턴

Utility 프로세스:

- CPU 집약적 백그라운드 작업을 위한 utility 프로세스 생성
- MessageChannel을 통해 utility 프로세스와 통신
- utility 프로세스 수명 주기 및 에러 처리

네이티브 모듈 통합:

- Electron Node.js 버전에 맞게 네이티브 모듈 재빌드
- better-sqlite3로 로컬 데이터베이스 저장
- keytar로 보안 자격 증명 저장

Protocol Handlers 및 Deep Linking:

- 앱 시작을 위한 커스텀 URL 프로토콜 등록
- 다양한 플랫폼에서 deep link 처리
- 커스텀 프로토콜을 통한 OAuth 콜백

성능 최적화:

- 무거운 모듈 및 창의 lazy loading
- 시작 시간 개선을 위한 deferred 초기화
- 장시간 실행 애플리케이션을 위한 메모리 관리

---

## Works Well With

- do-lang-typescript - 타입 안전 Electron 개발을 위한 TypeScript 패턴
- do-domain-frontend - React, Vue, 또는 Svelte renderer 개발
- do-lang-javascript - Main process용 Node.js 패턴
- do-domain-backend - Backend API 통합
- do-workflow-testing - 데스크톱 앱 테스트 전략

---

## Troubleshooting

일반적인 문제:

시작 시 흰 화면:

preload 스크립트 경로가 빌드된 출력 디렉토리에 맞게 올바르게 구성되었는지 확인하세요. loadFile 또는 loadURL 경로가 기존 파일을 가리키는지 확인하세요. DevTools를 열어 console 에러를 검사하세요. script 실행을 차단할 수 있는 CSP 설정을 검토하세요.

IPC 작동 안 함:

Main handler와 Renderer invocation에서 채널 이름이 정확히 일치하는지 확인하세요. 윈도가 로드되기 전에 handler가 등록되었는지 확인하세요. contextBridge 사용이 올바른 패턴을 따르는지 확인하세요 (exposeInMainWorld 사용).

Native modules 실패:

npm install 후 npm rebuild를 실행하여 네이티브 모듈을 재빌드하세요. Electron에 포함된 Node.js 버전과 일치하도록 합니다. postinstall 스크립트로 자동 재빌드를 자동화하세요.

Auto-Update 작동 안 함:

애플리케이션 code sign되었는지 확인하세요 (업데이트가 필요함). electron-builder.yml에서 publish 구성을 확인하세요. 연결 문제를 진단하기 위해 electron-updater logging을 활성화하세요. 업데이트 확인을 차단할 수 있는 방화벽 설정을 검토하세요.

디버그 명령어:

npx electron-rebuild로 네이티브 모듈을 재빌드하세요. npx electron --version으로 Electron 버전을 확인하세요. DEBUG=electron-updater 환경변수로 상세 업데이트 로깅을 활성화하세요.

---

## Resources

완전한 코드 예제 및 구성 템플릿은 다음을 참조하세요:

- reference.md - 상세 API 문서, 버전 매트릭스, Context7 라이브러리 매핑
- examples.md - 모든 패턴에 대한 프로덕션 준비 코드 예제

최신 문서의 경우 Context7를 사용하여 쿼리하세요:

- /electron/electron core Electron API
- /electron/forge Electron Forge 툴링
- /electron-userland/electron-builder 패키징 구성

---

Version: 2.0.0
Last Updated: 2026-01-10
Changes: CLAUDE.md 문서 표준을 준수하도록 재구성 - 모든 코드 예제 제거, 내러러러티 텍스트 형식으로 변환

---
name: do-lang-java
description: >
  Java 21 LTS 개발 전문가 - Spring Boot 3.3, virtual threads, pattern matching,
  엔터프라이즈 패턴涵盖. 엔터프라이즈 애플리케이션, 마이크로서비스,
  Spring 프로젝트 구축 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.1.0"
  category: "language"
  status: "active"
  updated: "2026-01-11"
  modularized: "false"
  tags: "java, spring-boot, jpa, hibernate, virtual-threads, enterprise"
  context7-libraries: "/spring-projects/spring-boot, /spring-projects/spring-framework, /spring-projects/spring-security"
  related-skills: "do-lang-kotlin, do-domain-backend"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["Java", "Spring Boot", "Spring Framework", "JPA", "Hibernate", "Maven", "Gradle", ".java", "pom.xml", "build.gradle", "virtual thread"]
  languages: ["java"]
---

## Quick Reference (30초 요약)

Spring Boot 3.3, Virtual Threads, 모던 Java 기능과 엔터프라이즈 개발을 위한 Java 21 LTS 전문가.

자동 트리거: .java 확장자 파일, 빌드 파일 (pom.xml, build.gradle, build.gradle.kts)

핵심 기능:

- Java 21 LTS: Virtual threads, pattern matching, record patterns, sealed classes
- Spring Boot 3.3: REST controller, 서비스, 리포지토리, WebFlux reactive
- Spring Security 6: JWT 인증, OAuth2, 역할 기반 액세스 제어
- JPA/Hibernate 7: 엔티티 매핑, 관계, 쿼리, 트랜잭션
- JUnit 5: 단위 테스트, 모의, TestContainers 통합
- 빌드 도구: Maven 3.9, Gradle 8.5 Kotlin DSL

---

## Implementation Guide (5분 가이드)

### Java 21 LTS 기능

Project Loom과 함께 Virtual Threads:

Executors.newVirtualThreadPerTaskExecutor에서 try-with-resources를 사용하세요. 0에서 10000까지 IntStream.range를 호출하고 forEach를 사용하여 1초 동안 sleep하고 반복 값을 반환하는 태스크를 제출하세요.

Structured Concurrency Preview 패턴:

new StructuredTaskScope.ShutdownOnFailure에서 try-with-resources를 사용하세요. lambda 식으로 fetchUser와 fetchOrders를 위해 scope.fork를 호출하여 태스크를 포크하세요. scope.join을 호출한 다음 throwIfFailed를 호출하세요. 두 태스크 공급자의 결과를 포함하는 새로운 복합 객체를 반환하세요.

Switch를 위한 Pattern Matching:

Object 파라미터를 받는 describe 메서드를 생성하세요. i > 0인 가드 조건이 있는 Integer i에 대해 "양수 정수" 메시지를 반환하는 케이스, Integer i에 대해 "음수 또는 0" 메시지를 반환하는 케이스, String s에 대해 "길이" 메시지를 반환하는 케이스, 와일드카드가 있는 List에 대해 "크기" 메시지를 반환하는 케이스, null에 대해 "null 값"을 반환하는 케이스, 기본값에 대해 "알 수 없는 타입"을 반환하는 switch 식을 사용하세요.

Record Patterns 및 Sealed Classes:

int x와 int y를 가진 Point record를 정의하세요. Point topLeft와 Point bottomRight를 가진 Rectangle record를 정의하세요. 두 Point 구성 요소를 변수로 분해하는 Rectangle 패턴과 함께 switch를 사용하여 너비 × 높이의 절대값을 반환하는 area 메서드를 생성하세요. Circle과 Rectangle를 허용하는 sealed Shape 인터페이스를 정의하세요. PI × radius²를 사용하는 area 메서드로 Circle record를 구현하세요.

### Spring Boot 3.3

REST Controller 패턴:

@RestController와 @RequestMapping("api/users") 어노테이션과 함께 UserController를 생성하세요. @RequiredArgsConstructor로 UserService를 주입하세요. @GetMapping과 @PathVariable이 있는 getUser 메서드를 생성하고 findById 결과를 ok로 매핑하거나 notFound를 반환하는 ResponseEntity를 반환하세요. @Valid와 @RequestBody가 있는 CreateUserRequest를 받는 createUser 메서드를 생성하고 URI 위치를 빌드하고 본문과 함께 created 응답을 반환하세요. @DeleteMapping이 있는 deleteUser 메서드를 생성하고 서비스 결과에 따라 noContent 또는 notFound를 반환하세요.

Service Layer 패턴:

@Service, @RequiredArgsConstructor, @Transactional(readOnly = true) 어노테이션과 함께 UserService를 생성하세요. UserRepository와 PasswordEncoder를 주입하세요. Optional을 반환하는 findById 메서드를 생성하세요. @Transactional create 메서드를 생성하고 중복 이메일 확인 시 DuplicateEmailException을 발생시키고, builder 패턴으로 비밀번호를 인코딩하여 User를 빌드하고 리포지토리에 저장합니다.

### Spring Security 6

Security Configuration 패턴:

@Configuration과 @EnableWebSecurity 어노테이션이 있는 SecurityConfig를 생성하세요. HttpSecurity를 받는 filterChain Bean을 정의하세요. 공개 API 경로에 permitAll, admin 경로에 hasRole("ADMIN"), 다른 모든 요청에 authenticated로 authorizeHttpRequests를 구성하세요. oauth2ResourceServer에 jwt default를 구성하세요. STATELESS로 sessionManagement를 설정하고 csrf를 비활성화하세요. BCryptPasswordEncoder를 반환하는 passwordEncoder Bean을 정의하세요.

### JPA/Hibernate 패턴

Entity Definition 패턴:

@Entity와 @Table 어노테이션이 있는 User 엔티티를 생성하세요. Lombok @Getter, @Setter, @NoArgsConstructor, @Builder 어노테이션을 추가하세요. @Id와 @GeneratedValue(IDENT)로 id를 정의하세요. @Column(nullable = false)로 name과 email을 정의하고 email도 unique로 설정하세요. @Enumerated(STRING)으로 status를 정의하세요. @OneToMany(mappedBy = "user", cascade = ALL, orphanRemoval = true)로 orders를 정의하세요.

Custom Queries와 함께 Repository:

JpaRepository를 확장하는 UserRepository를 생성하세요. Optional을 반환하는 findByEmail을 추가하세요. boolean을 반환하는 existsByEmail을 추가하세요. @Param과 JPQL LEFT JOIN FETCH가 있는 @Query 어노테이션으로 findByIdWithOrders를 추가하세요. Pageable을 받는 findByNameContainingIgnoreCase를 추가하세요.

Records로 DTO:

id, name, email, status를 가진 UserDto record를 생성하세요. User 엔티티에서 record를 구성하는 from 정적 팩토리 메서드를 추가하세요. @NotBlank와 @Size가 있는 name, @Email와 @NotBlank가 있는 email, @NotBlank와 @Size(min = 8)가 있는 password를 가진 CreateUserRequest record를 생성하세요.

---

## Advanced Patterns

### Virtual Threads 통합

@Service와 @RequiredArgsConstructor 어노테이션이 있는 AsyncUserService를 생성하세요. StructuredTaskScope.ShutdownOnFailure를 사용하여 fetchUserDetails 메서드를 생성하세요. user와 orders 쿼리를 위해 태스크를 포크하고, 실패 시 throw하고, 복합 결과를 반환하세요. newVirtualThreadPerTaskExecutor를 사용하여 processUsersInParallel 메서드를 생성하고 user ID를 스트리밍하여 처리 태스크를 제출합니다.

### Build Configuration

Maven 3.9 패턴:

spring-boot-starter-parent 버전 3.3.0을 parent로 하는 프로젝트를 정의하세요. java.version 속성을 21로 설정하세요. spring-boot-starter-web과 spring-boot-starter-data-jpa 의존성을 추가하세요.

Gradle 8.5 Kotlin DSL 패턴:

org.springframework.boot, io.spring.dependency-management, java 플러그인을 적용하세요. toolchain languageVersion을 21로 설정하세요. web 및 data-jpa 스타터를 위한 implementation 의존성, test 스타터를 위한 testImplementation 의존성을 추가하세요.

### JUnit 5와 함께 테스트

Unit Testing 패턴:

@ExtendWith(MockitoExtension)가 있는 테스트 클래스를 생성하세요. @Mock로 UserRepository를, @InjectMocks로 UserService를 추가하세요. existsByEmail이 false를 반환하고 save가 id가 있는 user를 반환하도록 stub한 shouldCreateUser 테스트를 생성하세요. service create를 호출하고 assertThat result.id가 1과 같은지 확인하세요.

TestContainers와 함께 Integration Testing:

@Testcontainers와 @SpringBootTest 어노테이션이 있는 테스트 클래스를 생성하세요. postgres:16-alpine 이미지로 static Container<PostgreSQL>를 정의하세요. datasource.url을 컨테이너에서 설정하도록 @DynamicPropertySource를 추가하세요. repository를 autowire하세요. user를 저장하고 id가 notNull인지 assert하는 테스트를 생성하세요.

---

## Context7 통합

최신 문서를 위한 라이브러리 매핑:

- spring-projects/spring-boot Spring Boot 3.3 문서
- spring-projects/spring-framework Spring Framework 코어
- spring-projects/spring-security Spring Security 6
- hibernate/hibernate-orm Hibernate 7 ORM 패턴
- junit-team/junit5 JUnit 5 테스팅 프레임워크

---

## Works Well With

- do-lang-kotlin Kotlin 상호 운용성 및 Spring Kotlin 확장
- do-domain-backend REST API, GraphQL, 마이크로서비스 아키텍처
- do-domain-database JPA, Hibernate, R2DBC 패턴
- do-foundation-quality JUnit 5, Mockito, TestContainers 통합
- do-infra-docker JVM 컨테이너 최적화

---

## Troubleshooting

일반적인 문제:

- 버전 불일치: java -version을 실행하고 JAVA_HOME이 Java 21을 가리키는지 확인하세요
- 컴파일 에러: mvn clean compile -X 또는 gradle build --info를 실행하세요
- Virtual thread 문제: 필요한 경우 Java 21+와 --enable-preview로 확인하세요
- JPA lazy loading: @Transactional 어노테이션 또는 JOIN FETCH 쿼리를 사용하세요

성능 팁:

- spring.threads.virtual.enabled를 true로 설정하여 Virtual Threads를 활성화하세요
- 더 빠른 시작을 위해 GraalVM Native Image를 사용하세요
- HikariCP로 연결 풀을 구성하세요

---

## Advanced Documentation

포괄적인 참조 자료:

- reference.md Java 21 기능, Context7 매핑, 성능
- examples.md 프로덕션 준비 Spring Boot 예제

---

Last Updated: 2026-01-11
Status: Production Ready (v1.1.0)

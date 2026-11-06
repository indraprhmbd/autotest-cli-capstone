# AutoTest-CLI: AI Agent Generator Unit Test dan Edge-Case Otomatis (Go)

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Agent Runtime](https://img.shields.io/badge/IBM_Bob-CLI_%26_ACP-052FAD?style=flat-square&logo=ibm&logoColor=white)
![Editor](https://img.shields.io/badge/Zed_Editor-ACP_Ready-084CCF?style=flat-square)
![Track](https://img.shields.io/badge/Track-IT_AI_Agent_for_Programming-2EA44F?style=flat-square)
![Coverage](https://img.shields.io/badge/Test_Coverage-100%25-brightgreen?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)

**Capstone Project: Hacktiv8 x IBM SkillsBuild**
* **Jenis Course:** IT - AI Agent for Programming
* **Use Case Topic:** Test Case Generator
* **Tools:** IBM Bob CLI (Bob Shell) via Agent Client Protocol (ACP)

---

## Ringkasan Project

**AutoTest-CLI** adalah AI Agent berbasis command-line dan Agent Client Protocol (ACP) yang dikembangkan menggunakan **IBM Bob CLI** untuk mengotomatisasi pembuatan unit test (`_test.go`) serta skenario pengujian edge-case langsung dari kode sumber berbahasa Go. Proyek ini bertujuan untuk mempercepat siklus pengujian perangkat lunak, meningkatkan cakupan pengujian (test coverage), dan meminimalkan potensi bug yang lolos ke tahap produksi akibat kurangnya pengujian pada kondisi batas (boundary/edge cases).

Proyek ini dirancang untuk digunakan oleh Software Engineer, Backend Developer, dan QA Automation Engineer yang membutuhkan pembuatan test suite secara cepat dan terstandarisasi langsung dari terminal atau editor berbasis ACP seperti Zed. Fungsi utama dari AI Agent ini meliputi analisis statis terhadap fungsi dan alur logika pada kode sumber Go, identifikasi potensi titik kegagalan (input kosong, batas karakter maksimal, isolasi tenant/room, dan error handling), serta pembuatan kode table-driven unit test yang langsung dieksekusi dan diverifikasi secara otomatis menggunakan `go test`.

## Problem Statement

Dalam siklus pengembangan perangkat lunak yang cepat, penulisan unit test secara manual sering kali terabaikan atau dilakukan secara tidak menyeluruh karena keterbatasan waktu developer. Kondisi saat ini menunjukkan bahwa mayoritas pengujian yang ditulis secara manual hanya berfokus pada skenario normal (happy path), sementara skenario batas (edge cases), penanganan error (negative path), dan input anomali sering terlewatkan. Selain itu, menulis boilerplate pengujian secara berulang menguras waktu produktif engineer.

Masalah ini sangat penting untuk diselesaikan karena ketiadaan pengujian yang komprehensif membuat sistem, terutama layanan middleware dan router AI Agent, rentan terhadap regression bug setiap kali ada perubahan kode baru. Apabila tidak ditangani, hal ini berdampak langsung pada meningkatnya risiko kegagalan sistem di lingkungan produksi (downtime), potensi kebocoran akses antar-tenant, tingginya biaya perbaikan bug pada tahap akhir, serta menurunnya keandalan aplikasi secara keseluruhan.

## Target Users

1. **Software / Backend Developer:** Pengembang yang menulis logika bisnis inti dalam bahasa Go dan membutuhkan asisten otomatis di terminal atau editor untuk menghasilkan unit test berstandar tinggi segera setelah selesai menulis fungsi baru.
2. **QA Automation Engineer:** Penguji perangkat lunak yang membutuhkan referensi skenario edge-case dan negative test secara menyeluruh berdasarkan struktur logika kode aktual untuk meningkatkan code coverage.
3. **Solo Developer / Systems Architect:** Pengembang independen yang mengelola infrastruktur microservice dan perlu menjaga kualitas serta stabilitas codebase secara efisien tanpa tim QA khusus.

---

## Arsitektur Prompt (Prompt Layering)

Arsitektur prompt pada **AutoTest-CLI** menerapkan pendekatan **4-Layer Prompting** yang didefinisikan secara deklaratif melalui file `AGENTS.md` dan dieksekusi melalui IBM Bob CLI:

| Layer | Nama Layer | Implementasi Prompt | Tujuan |
| :--- | :--- | :--- | :--- |
| **Layer 1** | **Role & Persona** | *"You are a Senior Go QA & Reliability Engineer specializing in deterministic unit testing, table-driven tests, and edge-case discovery for lightweight AI agent services."* | Menetapkan identitas dan keahlian spesifik AI Agent agar menghasilkan kode pengujian yang mengikuti standar idiomatic Go dan deterministik. |
| **Layer 2** | **Technical Constraints** | *"Use ONLY Go's standard `testing` package. Do NOT import external assertion libraries. Enforce idiomatic Go Table-Driven Tests (`[]struct{ name string; ... }`) executed with `t.Run()`. Never modify the target source file."* | Membatasi ruang gerak agen agar tidak menggunakan library pihak ketiga di luar `go.mod`, memastikan konsistensi struktur table-driven test, dan menjaga keutuhan kode sumber asli. |
| **Layer 3** | **Analytical & Edge-Case Matrix** | *"Before generating the test file, analyze the target Go source file to cover: (1) Happy path, (2) Boundary conditions (whitespace strings, `MaxPayloadLen <= 0` fallback, exact max length vs `maxLen + 1`), and (3) Negative/security cases (empty tenant, unauthorized room)."* | Memandu AI Agent memetakan seluruh cabang logika terlebih dahulu sehingga skenario pengujian mencakup alur normal maupun kondisi batas. |
| **Layer 4** | **Execution & Verification** | *"Generate `validator_test.go` in the same package, then execute `go test -v -cover ./...` via the CLI tool to verify that all generated test cases compile, pass, and achieve >95% statement coverage."* | Menginstruksikan IBM Bob CLI untuk menuliskan kode ke sistem file sekaligus memvalidasi kebenaran sintaks dan logika pengujian secara langsung di terminal. |

## Iterasi Prompt

Selama proses pengembangan di IBM Bob CLI, dilakukan 3 tahap iterasi untuk menyempurnakan kualitas kode pengujian yang dihasilkan:

1. **Iterasi 1 (Baseline Prompt)**
   * **Prompt Awal:** *"Buatkan unit test untuk file `validator.go` ini menggunakan bahasa Go."*
   * **Alasan Perubahan:** Output yang dihasilkan masih menulis fungsi `TestXxx` secara terpisah-pisah (repetitif), hanya menguji skenario sukses (happy path), dan mengimpor library eksternal `stretchr/testify` yang menyebabkan error saat kompilasi karena tidak terdaftar di `go.mod`.
2. **Iterasi 2 (Penambahan Constraint Standar Go dan Edge-Case)**
   * **Perubahan Prompt:** Menambahkan aturan tegas untuk hanya menggunakan standard library `testing`, mewajibkan pola Table-Driven Test (`t.Run`), serta meminta cakupan negative case (input kosong, ruangan tidak terotorisasi, batas kuota token).
   * **Alasan Perubahan:** Menyamakan standar kode dengan best practice Go modern dan menghilangkan ketergantungan pada library eksternal.
   * **Dampak:** Kode `validator_test.go` berhasil dikompilasi tanpa dependensi tambahan, struktur kode jauh lebih bersih, dan mencakup skenario input tidak valid. Namun, agen belum memverifikasi apakah ekspektasi error pada tes sudah sesuai dengan perilaku aktual kode secara mandiri.
3. **Iterasi 3 (Final: Integrasi Self-Verification Loop via CLI)**
   * **Perubahan Prompt:** Menambahkan Layer 4 pada `AGENTS.md` yang mewajibkan IBM Bob CLI menjalankan perintah `go test -v -cover ./...` setelah membuat file `_test.go` dan melakukan perbaikan mandiri apabila ditemukan kegagalan assertion.
   * **Alasan Perubahan:** Memastikan output dari AI Agent benar-benar siap pakai (production-ready) tanpa intervensi perbaikan manual oleh developer.
   * **Dampak:** Seluruh skenario happy path, boundary, dan negative case langsung berstatus `PASS` saat dieksekusi dengan tingkat cakupan pengujian (statement coverage) mencapai **100.0%**.

---

## Struktur Repository dan Cara Penggunaan

```text
autotest-cli-capstone/
├── AGENTS.md           # Konfigurasi 4-Layer Prompt Architecture untuk IBM Bob CLI
├── go.mod              # Definisi modul Go
├── validator.go        # Kode sumber target (AI Agent Tenant & Room Dispatch Validator)
├── validator_test.go   # Unit test hasil generate otomatis oleh IBM Bob CLI
└── README.md           # Dokumentasi Capstone Project
```

## Integrasi dengan Zed Editor (Agent Client Protocol)

Proyek ini dijalankan menggunakan IBM Bob CLI yang diintegrasikan ke Zed Editor melalui Agent Client Protocol (ACP):

* **Command:** `bob`
* **Arguments:** `acp --trust`

## Eksekusi Pengujian Manual

Untuk memverifikasi hasil unit test yang telah dibuat oleh IBM Bob CLI, jalankan perintah berikut di terminal:

```bash
go test -v -cover ./...
```

## Potensi Pengembangan Lanjutan

1. **Integrasi CI/CD Pre-Commit & GitHub Actions:** Mengintegrasikan AutoTest-CLI ke dalam pipeline otomatis untuk men-generate dan memvalidasi unit test setiap kali terdapat penambahan fungsi baru pada Pull Request.

2. **Automated Native Fuzz Testing (`go test -fuzz`):** Memperluas layer prompt agar AI Agent turut men-generate fungsi `FuzzXxx` bawaan Go 1.18+ guna menguji ketahanan parser dan validator terhadap mutasi payload acak secara kontinu.

3. **Mock Interface Generator Otomatis:** Menambahkan instruksi deteksi interface Go (seperti koneksi SQLite atau HTTP client) untuk men-generate lightweight mock struct secara otomatis tanpa memerlukan framework mocking eksternal.

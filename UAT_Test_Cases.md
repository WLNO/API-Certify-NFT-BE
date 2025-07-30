# User Acceptance Testing (UAT) - Certify-NFT Application

## Test Case Template
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|

---

## 1. Authentication Module

### TC-AUTH-001: User Registration
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-AUTH-001 | Authentication | User dapat mendaftar akun baru | 1. Buka halaman Register<br>2. Isi form dengan data valid<br>3. Klik "Register" | User berhasil terdaftar dan redirect ke login | User berhasil terdaftar | Passed | High | Gusti Padaka (BE Dev) |

### TC-AUTH-002: User Login
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-AUTH-002 | Authentication | User dapat login dengan kredensial valid | 1. Buka halaman Login<br>2. Masukkan email dan password<br>3. Klik "Login" | User berhasil login dan redirect ke dashboard | User berhasil login | Passed | High | Gusti Padaka (BE Dev) |

### TC-AUTH-003: Vendor Registration
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-AUTH-003 | Authentication | Vendor dapat mendaftar akun baru | 1. Buka halaman Register Vendor<br>2. Isi form vendor dengan data valid<br>3. Klik "Register" | Vendor berhasil terdaftar | Vendor berhasil terdaftar | Passed | High | Gusti Padaka (BE Dev) |

### TC-AUTH-004: Vendor Login
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-AUTH-004 | Authentication | Vendor dapat login ke dashboard vendor | 1. Buka halaman Vendor Login<br>2. Masukkan kredensial vendor<br>3. Klik "Login" | Vendor berhasil login ke dashboard vendor | Vendor berhasil login | Passed | High | Gusti Padaka (BE Dev) |

---

## 2. Event Management Module

### TC-EVENT-001: Create Event
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-EVENT-001 | Event Management | Vendor dapat membuat event baru | 1. Login sebagai vendor<br>2. Buka halaman Create Event<br>3. Isi semua field required<br>4. Upload gambar event<br>5. Klik "Create Event" | Event berhasil dibuat dan tersimpan | Event berhasil dibuat | Passed | High | Gusti Padaka (BE Dev) |

### TC-EVENT-002: View Event List
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-EVENT-002 | Event Management | User dapat melihat daftar event yang tersedia | 1. Login sebagai user<br>2. Buka halaman Events<br>3. Scroll daftar event | Daftar event ditampilkan dengan informasi lengkap | Daftar event ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-EVENT-003: Event Detail View
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-EVENT-003 | Event Management | User dapat melihat detail event | 1. Buka halaman Events<br>2. Klik pada event tertentu<br>3. Lihat detail event | Detail event ditampilkan dengan agenda, requirements, dan info lengkap | Detail event ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-EVENT-004: Manage Event
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-EVENT-004 | Event Management | Vendor dapat mengelola event yang dibuat | 1. Login sebagai vendor<br>2. Buka Manage Event<br>3. Lihat statistik dan kontrol panel | Dashboard manajemen event ditampilkan dengan fitur lengkap | Dashboard manajemen ditampilkan | Passed | High | Gusti Padaka (BE Dev) |

---

## 3. Certificate Management Module

### TC-CERT-001: Upload Certificates
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-CERT-001 | Certificate Management | Vendor dapat upload sertifikat untuk event | 1. Login sebagai vendor<br>2. Buka Manage Event<br>3. Pilih Upload Certificate<br>4. Upload file sertifikat<br>5. Klik "Upload" | Sertifikat berhasil diupload ke IPFS | Sertifikat berhasil diupload | Passed | High | Gusti Padaka (BE Dev) |

### TC-CERT-002: Mint Certificate NFT
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-CERT-002 | Certificate Management | Vendor dapat mint sertifikat menjadi NFT | 1. Setelah upload sertifikat<br>2. Klik "Mint Certificate"<br>3. Konfirmasi transaksi blockchain | NFT sertifikat berhasil di-mint ke blockchain | NFT berhasil di-mint | Passed | High | Gusti Padaka (BE Dev) |

### TC-CERT-003: User Certificate Minting
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-CERT-003 | Certificate Management | User dapat mint sertifikat ke wallet | 1. Login sebagai user<br>2. Buka My Certificates<br>3. Connect wallet<br>4. Klik "Mint Certificate"<br>5. Konfirmasi transaksi | Sertifikat NFT berhasil di-mint ke wallet user | Sertifikat berhasil di-mint ke wallet | Passed | High | Gusti Padaka (BE Dev) |

### TC-CERT-004: View My Certificates
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-CERT-004 | Certificate Management | User dapat melihat sertifikat yang dimiliki | 1. Login sebagai user<br>2. Buka My Certificates<br>3. Lihat daftar sertifikat | Daftar sertifikat ditampilkan dengan status minting | Daftar sertifikat ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

---

## 4. Certificate Verification Module

### TC-VERIFY-001: Verify Certificate
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-VERIFY-001 | Certificate Verification | User dapat verifikasi keaslian sertifikat | 1. Buka halaman Verify Certificate<br>2. Masukkan Certificate ID/Hash<br>3. Klik "Verify" | Hasil verifikasi ditampilkan dengan detail sertifikat | Verifikasi berhasil | Passed | High | Gusti Padaka (BE Dev) |

### TC-VERIFY-002: Verify Invalid Certificate
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-VERIFY-002 | Certificate Verification | Sistem dapat mendeteksi sertifikat palsu | 1. Buka halaman Verify Certificate<br>2. Masukkan ID sertifikat palsu<br>3. Klik "Verify" | Pesan error ditampilkan "Certificate not found" | Error message ditampilkan | Passed | Medium | Gusti Padaka (BE Dev) |

---

## 5. Whitelist Management Module

### TC-WHITELIST-001: Add to Whitelist
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-WHITELIST-001 | Whitelist Management | Vendor dapat menambah user ke whitelist | 1. Login sebagai vendor<br>2. Buka Manage Event<br>3. Pilih Whitelist Management<br>4. Tambah email user<br>5. Klik "Add to Whitelist" | User berhasil ditambahkan ke whitelist | User berhasil ditambahkan | Passed | High | Gusti Padaka (BE Dev) |

### TC-WHITELIST-002: View Whitelist
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-WHITELIST-002 | Whitelist Management | Vendor dapat melihat daftar whitelist | 1. Buka Whitelist Management<br>2. Lihat tabel whitelist | Daftar user dalam whitelist ditampilkan | Daftar whitelist ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-WHITELIST-003: Remove from Whitelist
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-WHITELIST-003 | Whitelist Management | Vendor dapat menghapus user dari whitelist | 1. Buka Whitelist Management<br>2. Klik "Remove" pada user tertentu<br>3. Konfirmasi penghapusan | User berhasil dihapus dari whitelist | User berhasil dihapus | Passed | Medium | Gusti Padaka (BE Dev) |

---

## 6. User Profile Module

### TC-PROFILE-001: View Profile
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-PROFILE-001 | User Profile | User dapat melihat profil sendiri | 1. Login sebagai user<br>2. Buka Profile page<br>3. Lihat informasi profil | Informasi profil user ditampilkan | Profil ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-PROFILE-002: Update Profile
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-PROFILE-002 | User Profile | User dapat mengupdate informasi profil | 1. Buka Profile page<br>2. Edit informasi profil<br>3. Klik "Save Changes" | Profil berhasil diupdate | Profil berhasil diupdate | Passed | Medium | Gusti Padaka (BE Dev) |

---

## 7. Dashboard Module

### TC-DASHBOARD-001: User Dashboard
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-DASHBOARD-001 | Dashboard | User dapat melihat dashboard pribadi | 1. Login sebagai user<br>2. Buka User Dashboard<br>3. Lihat statistik dan menu | Dashboard user ditampilkan dengan menu lengkap | Dashboard ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-DASHBOARD-002: Vendor Dashboard
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-DASHBOARD-002 | Dashboard | Vendor dapat melihat dashboard vendor | 1. Login sebagai vendor<br>2. Buka Vendor Dashboard<br>3. Lihat statistik event dan sertifikat | Dashboard vendor ditampilkan dengan fitur manajemen | Dashboard vendor ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

---

## 8. Wallet Integration Module

### TC-WALLET-001: Connect Wallet
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-WALLET-001 | Wallet Integration | User dapat connect wallet crypto | 1. Buka halaman yang memerlukan wallet<br>2. Klik "Connect Wallet"<br>3. Pilih MetaMask<br>4. Approve connection | Wallet berhasil terhubung | Wallet berhasil terhubung | Passed | High | Gusti Padaka (BE Dev) |

### TC-WALLET-002: View Wallet Info
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-WALLET-002 | Wallet Integration | User dapat melihat informasi wallet | 1. Connect wallet<br>2. Lihat Wallet Info Card | Informasi wallet ditampilkan (address, balance) | Info wallet ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

---

## 9. Navigation Module

### TC-NAV-001: Navigation Menu
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-NAV-001 | Navigation | User dapat navigasi antar halaman | 1. Login sebagai user<br>2. Klik menu navigasi<br>3. Pindah antar halaman | Navigasi berfungsi dengan baik | Navigasi berfungsi | Passed | Low | Rama, Faruq, Mumtaz (FE Dev) |

### TC-NAV-002: Responsive Design
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-NAV-002 | Navigation | Aplikasi responsive di berbagai device | 1. Buka aplikasi di desktop<br>2. Buka di tablet<br>3. Buka di mobile | Layout menyesuaikan dengan ukuran layar | Layout responsive | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

---

## 10. Error Handling Module

### TC-ERROR-001: Invalid Login
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-ERROR-001 | Error Handling | Sistem menampilkan error untuk login invalid | 1. Masukkan email/password salah<br>2. Klik "Login" | Pesan error ditampilkan | Error message ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

### TC-ERROR-002: Network Error
| Test Case ID | Module | Test Scenario | Test Steps | Expected Result | Actual Result | Status | Priority | Tester |
|--------------|--------|---------------|------------|-----------------|---------------|--------|----------|--------|
| TC-ERROR-002 | Error Handling | Sistem menangani error jaringan | 1. Matikan internet<br>2. Coba akses fitur yang memerlukan API | Pesan error jaringan ditampilkan | Network error ditampilkan | Passed | Medium | Rama, Faruq, Mumtaz (FE Dev) |

---

## Test Execution Summary

| Module | Total Test Cases | Passed | Failed | Not Tested | Pass Rate |
|--------|-----------------|--------|--------|------------|-----------|
| Authentication | 4 | 4 | 0 | 0 | 100% |
| Event Management | 4 | 4 | 0 | 0 | 100% |
| Certificate Management | 4 | 4 | 0 | 0 | 100% |
| Certificate Verification | 2 | 2 | 0 | 0 | 100% |
| Whitelist Management | 3 | 3 | 0 | 0 | 100% |
| User Profile | 2 | 2 | 0 | 0 | 100% |
| Dashboard | 2 | 2 | 0 | 0 | 100% |
| Wallet Integration | 2 | 2 | 0 | 0 | 100% |
| Navigation | 2 | 2 | 0 | 0 | 100% |
| Error Handling | 2 | 2 | 0 | 0 | 100% |
| **TOTAL** | **27** | **27** | **0** | **0** | **100%** |

---

## Test Environment Requirements

### Browser Compatibility
- Chrome (Latest)
- Firefox (Latest)
- Safari (Latest)
- Edge (Latest)

### Device Testing
- Desktop (1920x1080)
- Tablet (768x1024)
- Mobile (375x667)

### Network Conditions
- High-speed internet
- Slow internet (3G simulation)
- Offline mode

### Wallet Requirements
- MetaMask extension
- Test network (Goerli/Sepolia)
- Test ETH for gas fees

---

## Notes
- Semua test case dibuat berdasarkan fitur yang ada di codebase
- Priority: High (Critical), Medium (Important), Low (Nice to have)
- Status: Not Tested, Passed, Failed, Blocked
- Test cases dapat diupdate sesuai dengan perubahan requirement 
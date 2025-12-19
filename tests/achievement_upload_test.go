package tests

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"prestasi-mahasiswa-api/controllers"
	"prestasi-mahasiswa-api/models"
	"prestasi-mahasiswa-api/utils"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK DEFINITION ---
// Mendefinisikan MockAchieveService di dalam file yang sama untuk menghindari error undefined
type MockAchieveService struct {
	mock.Mock
}

func (m *MockAchieveService) AddAttachment(ctx context.Context, sid uuid.UUID, rid uuid.UUID, att models.AttachmentFile) (int, error) {
	args := m.Called(sid, rid, att)
	return args.Int(0), args.Error(1)
}

// Implementasi placeholder agar memenuhi interface AchievementService
func (m *MockAchieveService) CreateDraft(ctx context.Context, sid uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) { return nil, 0, nil }
func (m *MockAchieveService) UpdateDraft(ctx context.Context, sid uuid.UUID, rid uuid.UUID, req *models.CreateAchievementRequest) (*models.AchievementReference, int, error) { return nil, 0, nil }
func (m *MockAchieveService) DeleteDraft(ctx context.Context, sid uuid.UUID, rid uuid.UUID) (int, error) { return 0, nil }
func (m *MockAchieveService) SubmitForVerification(ctx context.Context, sid uuid.UUID, rid uuid.UUID) (*models.AchievementReference, int, error) { return nil, 0, nil }
func (m *MockAchieveService) ListFilteredAchievements(ctx context.Context, c *utils.JWTCustomClaims) ([]models.AchievementDetailResponse, int, error) { return nil, 0, nil }
func (m *MockAchieveService) GetDetailWithVerification(ctx context.Context, c *utils.JWTCustomClaims, rid uuid.UUID) (*models.AchievementDetailResponse, int, error) { return nil, 0, nil }
func (m *MockAchieveService) VerifyAchievement(ctx context.Context, aid uuid.UUID, rid uuid.UUID) (*models.AchievementReference, int, error) { return nil, 0, nil }
func (m *MockAchieveService) RejectAchievement(ctx context.Context, aid uuid.UUID, rid uuid.UUID, n string) (*models.AchievementReference, int, error) { return nil, 0, nil }
func (m *MockAchieveService) HardDelete(ctx context.Context, rid uuid.UUID) (int, error) { return 0, nil }

// --- TEST CASE ---
func TestUploadAttachmentFromTestsFolder(t *testing.T) {
	app := fiber.New()
	mockSvc := new(MockAchieveService)
	ctrl := controllers.AchievementController{Service: mockSvc}
	
	// Setup Route
	app.Post("/achievements/:id/attachments", ctrl.UploadAttachment)

	t.Run("Successfully Upload File from tests/attachment", func(t *testing.T) {
		// 1. Tentukan file sumber di tests/attachment
		// Catatan: Saat menjalankan 'go test ./tests/...', path relatif dimulai dari root project
		sourceFileName := "bukti_prestasi.jpg" // PASTIKAN FILE INI ADA DI tests/attachment/
		sourcePath := filepath.Join("attachment", sourceFileName)

		file, err := os.Open(sourcePath)
		if err != nil {
			t.Fatalf("File sumber tidak ditemukan di %s. Silakan taruh file %s di folder tests/attachment/", sourcePath, sourceFileName)
		}
		defer file.Close()

		// 2. Siapkan Form Data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("attachment", sourceFileName)
		assert.NoError(t, err)
		_, err = io.Copy(part, file)
		assert.NoError(t, err)
		writer.Close()

		// 3. Mock Expectations
		mockSvc.On("AddAttachment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(201, nil)

		// 4. Execute Request
		// Gunakan UUID dummy untuk testing
		testID := uuid.New().String()
		req := httptest.NewRequest("POST", "/achievements/"+testID+"/attachments", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		// Bypass middleware dengan manual context atau mock auth jika diperlukan
		// Di sini kita langsung panggil handler melalui app.Test
		resp, _ := app.Test(req)

		// 5. Verifikasi
		assert.Equal(t, 201, resp.StatusCode)
	})
}
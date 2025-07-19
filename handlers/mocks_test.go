package handlers

import (
	"competitions/models"
	"context"
	"errors"
)

// Mock do repositório de usuários
// Implementa todos os métodos da interface repository.UsuarioRepository
// Adapte os retornos conforme necessário para seus testes

type MockUsuarioRepository struct{}

func (m *MockUsuarioRepository) FindAll(ctx context.Context) ([]models.Usuario, error) {
	return []models.Usuario{{ID: 1, Nome: "Teste", Email: "teste@email.com"}}, nil
}
func (m *MockUsuarioRepository) FindByID(ctx context.Context, id int) (*models.Usuario, error) {
	return &models.Usuario{ID: 1, Nome: "Teste", Email: "teste@email.com"}, nil
}
func (m *MockUsuarioRepository) FindByIDForAuth(ctx context.Context, id uint) (*models.Usuario, error) {
	return &models.Usuario{ID: id, Nome: "Teste", Email: "teste@email.com"}, nil
}
func (m *MockUsuarioRepository) FindByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	if email == "teste@email.com" {
		return &models.Usuario{ID: 1, Email: email, Password: "$2a$10$b8/U.ay.jH.OfER4jD4.BOsPXBb.D5Rz.bp3i.jH.OfER4jD4.BO"}, nil
	}
	return nil, errors.New("usuário não encontrado")
}
func (m *MockUsuarioRepository) Create(ctx context.Context, usuario *models.Usuario) error {
	return nil
}
func (m *MockUsuarioRepository) Update(ctx context.Context, usuario *models.Usuario) (int64, error) {
	return 1, nil
}
func (m *MockUsuarioRepository) UpdatePassword(ctx context.Context, id uint, newPassword string) (int64, error) {
	return 1, nil
}
func (m *MockUsuarioRepository) Delete(ctx context.Context, id int) (int64, error) { return 1, nil }
func (m *MockUsuarioRepository) AssociateEsporte(ctx context.Context, usuarioID int, esporteIDs []int) error {
	return nil
}
func (m *MockUsuarioRepository) GetEsportesByUsuario(ctx context.Context, userID int) ([]models.Esporte, error) {
	return nil, nil
}
func (m *MockUsuarioRepository) GetUsuariosByEsporte(ctx context.Context, esporteID int) ([]models.Usuario, error) {
	return nil, nil
}

// Mock do repositório de torneios
type MockTorneioRepository struct{}

func (m *MockTorneioRepository) Create(ctx context.Context, input models.TorneioInput) (models.Torneio, error) {
	return models.Torneio{ID: 1, Nome: "Torneio Teste"}, nil
}
func (m *MockTorneioRepository) FindAll(ctx context.Context) ([]models.Torneio, error) {
	return nil, nil
}
func (m *MockTorneioRepository) FindByID(ctx context.Context, id int) (models.Torneio, error) {
	return models.Torneio{}, nil
}
func (m *MockTorneioRepository) Update(ctx context.Context, id int, input models.TorneioInput) (int64, error) {
	return 1, nil
}
func (m *MockTorneioRepository) Delete(ctx context.Context, id int) (int64, error) { return 1, nil }
func (m *MockTorneioRepository) InscreverJogador(ctx context.Context, inscricao models.JogadorTorneio) (models.JogadorTorneio, error) {
	return models.JogadorTorneio{}, nil
}
func (m *MockTorneioRepository) ListarInscricoesPorTorneio(ctx context.Context, torneioID int) ([]models.InscricaoDetalhada, error) {
	return nil, nil
}

// Mock do repositório de esportes
type MockEsporteRepository struct{}

func (m *MockEsporteRepository) FindAll(ctx context.Context) ([]models.Esporte, error) {
	return []models.Esporte{{ID: 1, Nome: "Futebol"}}, nil
}
func (m *MockEsporteRepository) Create(ctx context.Context, input models.EsporteInput) (*models.Esporte, error) {
	return &models.Esporte{ID: 1, Nome: input.Nome}, nil
}
func (m *MockEsporteRepository) FindByID(ctx context.Context, id int) (*models.Esporte, error) {
	return &models.Esporte{ID: id, Nome: "Futebol"}, nil
}
func (m *MockEsporteRepository) Update(ctx context.Context, id int, input models.EsporteInput) (int64, error) {
	return 1, nil
}
func (m *MockEsporteRepository) Delete(ctx context.Context, id int) (int64, error) { return 1, nil }

// Mock do repositório de grupos
type MockGrupoRepository struct{}

func (m *MockGrupoRepository) CreateGrupos(ctx context.Context, torneioID int, input models.CriarGruposInput) ([]models.GrupoComJogadores, error) {
	return []models.GrupoComJogadores{{Grupo: models.Grupo{ID: 1, Nome: "Grupo A"}}}, nil
}
func (m *MockGrupoRepository) GetEstatisticasGrupo(ctx context.Context, grupoID int) ([]models.EstatisticasJogador, error) {
	return nil, nil
}

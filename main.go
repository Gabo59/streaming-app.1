package main

import (
	"errors"
	"fmt"
	"strconv" // Necesario para convertir enteros a string para IDs
	"time"    // Para simular duraciones y marcas de tiempo
)

// =========================================================================
// 1. Definiciones de Tipos (Structs) y Encapsulación
//    Utilizamos campos no exportados (minúscula) para encapsular la data
//    y métodos para acceder a ella.
// =========================================================================

// Stream representa un elemento de contenido de streaming.
type Stream struct {
	id          string // id del stream (ej. "movie-1", "series-ep-5") - no exportado
	title       string // Título del contenido - no exportado
	genre       string // Género (ej. "Acción", "Comedia") - no exportado
	durationMin int    // Duración en minutos - no exportado
	url         string // URL de reproducción - no exportado
}

// NewStream es una función constructora para crear una nueva instancia de Stream.
// Esto promueve la encapsulación al controlar cómo se inicializan los objetos.
func NewStream(id, title, genre, url string, durationMin int) *Stream {
	return &Stream{
		id:          id,
		title:       title,
		genre:       genre,
		durationMin: durationMin,
		url:         url,
	}
}

// Métodos de acceso (Getters) para Stream, demostrando encapsulación.
func (s *Stream) GetID() string {
	return s.id
}
func (s *Stream) GetTitle() string {
	return s.title
}
func (s *Stream) GetGenre() string {
	return s.genre
}
func (s *Stream) GetDurationMin() int {
	return s.durationMin
}
func (s *Stream) GetURL() string {
	return s.url
}

// User representa un usuario del sistema de streaming.
type User struct {
	id            string   // ID del usuario - no exportado
	username      string   // Nombre de usuario - no exportado
	subscription  string   // Tipo de suscripción (ej. "Premium", "Basic") - no exportado
	watchHistory  []string // Slice de IDs de streams vistos - no exportado
	currentStream *Stream  // Stream actualmente en reproducción - no exportado
}

// NewUser es una función constructora para crear una nueva instancia de User.
func NewUser(id, username, subscription string) *User {
	return &User{
		id:           id,
		username:     username,
		subscription: subscription,
		watchHistory: []string{}, // Inicializa el historial como un slice vacío
	}
}

// Métodos de acceso (Getters) para User.
func (u *User) GetID() string {
	return u.id
}
func (u *User) GetUsername() string {
	return u.username
}
func (u *User) GetSubscription() string {
	return u.subscription
}

// AddToWatchHistory añade un stream al historial de visualización del usuario.
// Ejemplo de método que modifica el estado interno de forma controlada.
func (u *User) AddToWatchHistory(streamID string) {
	u.watchHistory = append(u.watchHistory, streamID)
}

// GetWatchHistory devuelve una copia del historial para mantener la encapsulación.
// Esto evita modificaciones externas directas del slice interno.
func (u *User) GetWatchHistory() []string {
	// Devuelve una copia para evitar que el slice interno sea modificado directamente desde fuera
	historyCopy := make([]string, len(u.watchHistory))
	copy(historyCopy, u.watchHistory)
	return historyCopy
}

// SetCurrentStream establece el stream actual que el usuario está viendo.
func (u *User) SetCurrentStream(s *Stream) {
	u.currentStream = s
}

// GetCurrentStream devuelve el stream actual que el usuario está viendo.
func (u *User) GetCurrentStream() *Stream {
	return u.currentStream
}

// =========================================================================
// 2. Manejo de Errores Personalizados (Implementando la interfaz error)
//    Esto permite errores más descriptivos y manejables.
// =========================================================================

var (
	ErrStreamNotFound = errors.New("stream not found")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidInput   = errors.New("invalid input data")
	ErrUnauthorized   = errors.New("unauthorized access")
)

// =========================================================================
// 3. Interfaces
//    Definimos interfaces para abstraer comportamientos.
// =========================================================================

// Playable define la capacidad de ser reproducido.
type Playable interface {
	GetTitle() string
	GetURL() string
	GetDurationMin() int
}

// Verify que Stream implementa Playable. Esto es opcional, pero útil.
var _ Playable = (*Stream)(nil)

// StreamStore define el contrato para almacenar y recuperar streams.
type StreamStore interface {
	AddStream(stream *Stream) error
	GetStreamByID(id string) (*Stream, error)
	GetAllStreams() []*Stream
}

// UserStore define el contrato para almacenar y recuperar usuarios.
type UserStore interface {
	AddUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetAllUsers() []*User
}

// =========================================================================
// 4. Funciones de Utilidad y Auxiliares (Principios funcionales: pureza)
//    Estas funciones son puras: no modifican el estado externo y siempre
//    devuelven la misma salida para la misma entrada.
// =========================================================================

// generateNextID genera un ID simple. En un sistema real, sería más robusto.
func generateNextID(prefix string, currentCount int) string {
	return prefix + "-" + strconv.Itoa(currentCount+1)
}

// isValidGenre simula una validación de género.
func isValidGenre(genre string) bool {
	validGenres := map[string]bool{
		"Accion":     true,
		"Comedia":    true,
		"Drama":      true,
		"Terror":     true,
		"Musical":    true,
		"Documental": true,
		"Animacion":  true,
	}
	return validGenres[genre]
}

// simulatePlaybackDuration simula la espera de reproducción de un stream.
func simulatePlaybackDuration(durationMin int) {
	fmt.Printf("Simulando reproducción por %d minutos...\n", durationMin)
	time.Sleep(time.Duration(durationMin) * time.Second) // Usamos segundos para simulación rápida
	fmt.Println("Reproducción finalizada.")
}

// =========================================================================
// 5. Lógica del Módulo de Contenido (`content` package / section)
//    Uso de maps para un acceso eficiente por ID.
// =========================================================================

// InMemoryStreamStore implementa StreamStore utilizando un map en memoria.
// Los campos son no exportados para encapsulación.
type InMemoryStreamStore struct {
	streams map[string]*Stream // map[ID del stream]Stream
	nextID  int                // Contador para generar IDs
}

// NewInMemoryStreamStore crea una nueva instancia del almacén de streams en memoria.
func NewInMemoryStreamStore() *InMemoryStreamStore {
	return &InMemoryStreamStore{
		streams: make(map[string]*Stream),
		nextID:  0,
	}
}

// AddStream añade un stream al almacén.
// Retorna un error si la entrada es inválida.
func (s *InMemoryStreamStore) AddStream(stream *Stream) error {
	if stream == nil || stream.GetTitle() == "" || stream.GetURL() == "" || stream.GetDurationMin() <= 0 {
		return ErrInvalidInput
	}
	if !isValidGenre(stream.GetGenre()) {
		return fmt.Errorf("%w: genero '%s' no valido", ErrInvalidInput, stream.GetGenre())
	}

	// Como el ID viene del constructor, verificamos si ya existe.
	// En este diseño, NewStream ya provee el ID, así que simplemente lo usamos.
	// Si quisiéramos auto-generar aquí, usaríamos generateNextID y lo asignaríamos al stream.
	if _, exists := s.streams[stream.GetID()]; exists {
		return fmt.Errorf("stream con ID '%s' ya existe", stream.GetID())
	}

	s.streams[stream.GetID()] = stream
	s.nextID++ // Incrementamos el contador para futuras auto-generaciones si se implementara
	return nil
}

// GetStreamByID recupera un stream por su ID.
// Retorna ErrStreamNotFound si el stream no existe.
func (s *InMemoryStreamStore) GetStreamByID(id string) (*Stream, error) {
	stream, ok := s.streams[id]
	if !ok {
		return nil, ErrStreamNotFound
	}
	return stream, nil
}

// GetAllStreams devuelve todos los streams almacenados.
// Retorna un slice de punteros a Stream.
func (s *InMemoryStreamStore) GetAllStreams() []*Stream {
	// Devuelve una copia del slice de streams para mantener la encapsulación.
	// No se modifica el map original directamente.
	allStreams := make([]*Stream, 0, len(s.streams))
	for _, stream := range s.streams {
		allStreams = append(allStreams, stream)
	}
	return allStreams
}

// =========================================================================
// 6. Lógica del Módulo de Usuario (`user` package / section)
//    Uso de maps para un acceso eficiente por ID de usuario.
// =========================================================================

// InMemoryUserStore implementa UserStore utilizando un map en memoria.
type InMemoryUserStore struct {
	users  map[string]*User // map[ID del usuario]User
	nextID int              // Contador para generar IDs
}

// NewInMemoryUserStore crea una nueva instancia del almacén de usuarios en memoria.
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:  make(map[string]*User),
		nextID: 0,
	}
}

// AddUser añade un usuario al almacén.
func (us *InMemoryUserStore) AddUser(user *User) error {
	if user == nil || user.GetUsername() == "" || user.GetSubscription() == "" {
		return ErrInvalidInput
	}
	if _, exists := us.users[user.GetID()]; exists {
		return fmt.Errorf("usuario con ID '%s' ya existe", user.GetID())
	}
	us.users[user.GetID()] = user
	us.nextID++
	return nil
}

// GetUserByID recupera un usuario por su ID.
func (us *InMemoryUserStore) GetUserByID(id string) (*User, error) {
	user, ok := us.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetAllUsers devuelve todos los usuarios registrados.
func (us *InMemoryUserStore) GetAllUsers() []*User {
	allUsers := make([]*User, 0, len(us.users))
	for _, user := range us.users {
		allUsers = append(allUsers, user)
	}
	return allUsers
}

// =========================================================================
// 7. Lógica del Módulo de Reproducción (`playback` package / section)
// =========================================================================

// PlayStream simula la reproducción de un stream para un usuario.
// Demuestra manejo de errores e interacción con objetos Stream y User.
func PlayStream(user *User, stream Playable) error {
	if user == nil {
		return ErrUserNotFound
	}
	if stream == nil {
		return ErrStreamNotFound
	}

	fmt.Printf("\n--- Iniciando reproducción para %s ---\n", user.GetUsername())
	fmt.Printf("Reproduciendo: %s (Género: %s, Duración: %d min)\n", stream.GetTitle(), stream.(*Stream).GetGenre(), stream.GetDurationMin()) // Casting para GetGenre
	fmt.Printf("URL: %s\n", stream.GetURL())

	user.SetCurrentStream(stream.(*Stream)) // Asignar el stream actual al usuario
	simulatePlaybackDuration(stream.GetDurationMin())
	user.AddToWatchHistory(stream.(*Stream).GetID()) // Añadir al historial
	user.SetCurrentStream(nil)                       // Limpiar stream actual al finalizar

	fmt.Printf("--- Reproducción de %s finalizada ---\n", stream.GetTitle())
	return nil
}

// =========================================================================
// 8. Estructura de Gestión del Sistema (`main` / `StreamingPlatform`)
//    Coordina los diferentes "módulos" (aquí representados por las stores).
// =========================================================================

// StreamingPlatform es la estructura principal que gestiona el sistema.
// Contiene instancias de los almacenes de datos.
type StreamingPlatform struct {
	streamStore StreamStore
	userStore   UserStore
}

// NewStreamingPlatform crea una nueva instancia de la plataforma de streaming.
func NewStreamingPlatform(ss StreamStore, us UserStore) *StreamingPlatform {
	return &StreamingPlatform{
		streamStore: ss,
		userStore:   us,
	}
}

// RegisterUser es una función de alto nivel para registrar un nuevo usuario.
func (p *StreamingPlatform) RegisterUser(username, subscription string) (*User, error) {
	newUserID := generateNextID("user", len(p.userStore.GetAllUsers()))
	newUser := NewUser(newUserID, username, subscription)
	err := p.userStore.AddUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("fallo al registrar usuario: %w", err)
	}
	fmt.Printf("Usuario registrado: %s (ID: %s, Suscripción: %s)\n", newUser.GetUsername(), newUser.GetID(), newUser.GetSubscription())
	return newUser, nil
}

// AddContent es una función de alto nivel para añadir nuevo contenido.
func (p *StreamingPlatform) AddContent(title, genre, url string, duration int) (*Stream, error) {
	newStreamID := generateNextID("stream", len(p.streamStore.GetAllStreams()))
	newStream := NewStream(newStreamID, title, genre, url, duration)
	err := p.streamStore.AddStream(newStream)
	if err != nil {
		return nil, fmt.Errorf("fallo al añadir contenido: %w", err)
	}
	fmt.Printf("Contenido añadido: %s (ID: %s, Duración: %d min)\n", newStream.GetTitle(), newStream.GetID(), newStream.GetDurationMin())
	return newStream, nil
}

// GetContentDetails es una función de alto nivel para obtener detalles de un stream.
func (p *StreamingPlatform) GetContentDetails(streamID string) (*Stream, error) {
	return p.streamStore.GetStreamByID(streamID)
}

// UserWatchStream simula la acción de un usuario viendo un stream.
func (p *StreamingPlatform) UserWatchStream(userID, streamID string) error {
	user, err := p.userStore.GetUserByID(userID)
	if err != nil {
		return err // Retorna ErrUserNotFound
	}
	stream, err := p.streamStore.GetStreamByID(streamID)
	if err != nil {
		return err // Retorna ErrStreamNotFound
	}

	return PlayStream(user, stream)
}

// =========================================================================
// 9. Función Principal (`main` function)
//    Aquí se orquesta la aplicación y se demuestra su uso.
// =========================================================================

func main() {
	fmt.Println("--- Sistema de Gestión de Streaming Iniciado ---")

	// 1. Inicializar los almacenes de datos
	streamStore := NewInMemoryStreamStore()
	userStore := NewInMemoryUserStore()

	// 2. Crear la plataforma de streaming
	platform := NewStreamingPlatform(streamStore, userStore)

	// --- DEMOSTRACIÓN DE FUNCIONALIDADES ---

	// a) Registrar Usuarios
	fmt.Println("\n--- Registro de Usuarios ---")
	user1, err := platform.RegisterUser("alice", "Premium")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	user2, err := platform.RegisterUser("bob", "Basic")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Intentar registrar un usuario con ID duplicado (fallará, pero el ID se genera secuencialmente aquí)
	// Para demostrar el manejo de errores, podríamos intentar añadir manualmente uno ya existente.
	// Por simplicidad, este ejemplo solo usa el generador secuencial.

	// b) Añadir Contenido
	fmt.Println("\n--- Añadiendo Contenido ---")
	movie1, err := platform.AddContent("Inception", "Accion", "http://stream.com/inception", 148)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	series1, err := platform.AddContent("Breaking Bad S1E1", "Drama", "http://stream.com/bb-s1e1", 55)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	_, err = platform.AddContent("Classical Mix", "Musical", "http://stream.com/classical", 60)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Intentar añadir contenido con género inválido
	_, err = platform.AddContent("Unknown Movie", "Fantasy", "http://stream.com/unknown", 90)
	if err != nil {
		fmt.Printf("Error: %v\n", err) // Debería mostrar ErrInvalidInput
	}

	// c) Listar Contenido Disponible
	fmt.Println("\n--- Contenido Disponible ---")
	for _, s := range streamStore.GetAllStreams() {
		fmt.Printf("ID: %s, Título: %s, Género: %s, Duración: %d min\n", s.GetID(), s.GetTitle(), s.GetGenre(), s.GetDurationMin())
	}

	// d) Simular Reproducción de Streams
	fmt.Println("\n--- Simulación de Reproducción ---")
	if user1 != nil && movie1 != nil {
		err = platform.UserWatchStream(user1.GetID(), movie1.GetID())
		if err != nil {
			fmt.Printf("Error al reproducir para %s: %v\n", user1.GetUsername(), err)
		}
	}

	if user2 != nil && series1 != nil {
		err = platform.UserWatchStream(user2.GetID(), series1.GetID())
		if err != nil {
			fmt.Printf("Error al reproducir para %s: %v\n", user2.GetUsername(), err)
		}
	}

	// Intentar reproducir un stream que no existe
	if user1 != nil {
		err = platform.UserWatchStream(user1.GetID(), "stream-999")
		if err != nil {
			fmt.Printf("Error al reproducir stream inexistente: %v\n", err) // Debería mostrar ErrStreamNotFound
		}
	}

	// Intentar que un usuario inexistente reproduzca algo
	if movie1 != nil {
		err = platform.UserWatchStream("user-999", movie1.GetID())
		if err != nil {
			fmt.Printf("Error al reproducir con usuario inexistente: %v\n", err) // Debería mostrar ErrUserNotFound
		}
	}

	// e) Ver Historial de Usuarios
	fmt.Println("\n--- Historial de Visualización ---")
	if user1 != nil {
		fmt.Printf("Historial de %s: %v\n", user1.GetUsername(), user1.GetWatchHistory())
	}
	if user2 != nil {
		fmt.Printf("Historial de %s: %v\n", user2.GetUsername(), user2.GetWatchHistory())
	}

	fmt.Println("\n--- Sistema de Gestión de Streaming Finalizado ---")
}

// Package labdeps pins the third-party modules that lab missions may import.
// Run `make deps` once (needs network) so labs can build offline.
package labdeps

import (
	_ "github.com/alicebob/miniredis/v2"
	_ "github.com/coder/websocket"
	_ "github.com/gin-gonic/gin"
	_ "github.com/glebarez/sqlite"
	_ "github.com/gofiber/fiber/v2"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/segmentio/kafka-go"
	_ "go.mongodb.org/mongo-driver/v2/mongo"
	_ "golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/ssh"
	_ "golang.org/x/image/draw"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/test/bufconn"
	_ "google.golang.org/protobuf/proto"
	_ "gorm.io/gorm"
	_ "modernc.org/sqlite"
)

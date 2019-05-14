package server

import (
	context "context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/mas9612/authserver/pkg/authserver"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/ldap.v3"
)

const (
	internalServerErrMsg = "internal server error occured"
)

// Authserver is the implementation of pb.AuthserverServer.
type Authserver struct {
	logger     *zap.Logger
	ldapaddr   string
	ldapport   int
	userFormat string
	pemPath    string
	pubkeyPath string
	issuer     string
	audience   string
}

type config struct {
	ldapaddr   string
	ldapport   int
	userFormat string
	pem        string
	pubkey     string
	issuer     string
	audience   string
}

// Option is the option to create authserver instance.
type Option func(*config)

// SetAddr sets the LDAP server address.
func SetAddr(addr string) Option {
	return func(c *config) {
		c.ldapaddr = addr
	}
}

// SetPort sets the LDAP server port.
func SetPort(port int) Option {
	return func(c *config) {
		c.ldapport = port
	}
}

// SetUserFormat sets the user format used when bind to LDAP server.
func SetUserFormat(format string) Option {
	return func(c *config) {
		c.userFormat = format
	}
}

// SetPem sets the pem path used to sign JWT token.
func SetPem(pem string) Option {
	return func(c *config) {
		c.pem = pem
	}
}

// SetPubkey sets the pubkey path used to validate JWT token.
// Default is authserver.pub located in same directory as authserver.
func SetPubkey(pubkey string) Option {
	return func(c *config) {
		c.pubkey = pubkey
	}
}

// SetIssuer sets the issuer used in JWT claim.
func SetIssuer(issuer string) Option {
	return func(c *config) {
		c.issuer = issuer
	}
}

// NewAuthserver creates new server instance.
func NewAuthserver(logger *zap.Logger, opts ...Option) (pb.AuthserverServer, error) {
	c := config{
		ldapaddr:   "localhost",
		ldapport:   389,
		userFormat: "%s",
		pubkey:     "authserver.pub",
	}
	for _, o := range opts {
		o(&c)
	}

	return &Authserver{
		logger:     logger,
		ldapaddr:   c.ldapaddr,
		ldapport:   c.ldapport,
		userFormat: c.userFormat,
		pemPath:    c.pem,
		pubkeyPath: c.pubkey,
	}, nil
}

// AuthClaim represents claim of JWT token.
type AuthClaim struct {
	User string `json:"user"`
	jwt.StandardClaims
}

// CreateToken creates and returns the new token.
func (s *Authserver) CreateToken(ctx context.Context, req *pb.CreateTokenRequest) (*pb.Token, error) {
	if req.User == "" || req.Password == "" || req.OrigHost == "" {
		errMsg := "too few argument: user, password, orig_host are required"
		s.logger.Error(errMsg)
		return nil, status.Error(codes.InvalidArgument, errMsg)
	}

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", s.ldapaddr, s.ldapport))
	if err != nil {
		errMsg := "failed to connect to LDAP server"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Internal, errMsg)
	}
	defer conn.Close()

	if err := conn.Bind(fmt.Sprintf(s.userFormat, req.User), req.Password); err != nil {
		errMsg := "bind failed"
		s.logger.Error(errMsg, zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, errMsg)
	}

	signKeyBytes, err := ioutil.ReadFile(s.pemPath)
	if err != nil {
		s.logger.Error("failed to load signkey", zap.Error(err))
		return nil, status.Error(codes.Internal, internalServerErrMsg)
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	if err != nil {
		s.logger.Error("failed to parse signkey", zap.Error(err))
		return nil, status.Error(codes.Internal, internalServerErrMsg)
	}

	nowUnix := time.Now().Unix()
	v4 := uuid.NewV4()
	claims := AuthClaim{
		req.User,
		jwt.StandardClaims{
			Audience:  req.OrigHost,
			ExpiresAt: nowUnix + 3600, // valid 1h
			Id:        v4.String(),
			IssuedAt:  nowUnix,
			Issuer:    s.issuer,
			NotBefore: nowUnix - 5,
			Subject:   "access_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		s.logger.Error("failed to generate JWT token", zap.Error(err))
		return nil, status.Error(codes.Internal, internalServerErrMsg)
	}

	return &pb.Token{
		Token: ss,
	}, nil
}

// ValidateToken validates given token and returns its validity.
func (s *Authserver) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claim := AuthClaim{}
	_, err := jwt.ParseWithClaims(in.Token, &claim, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("requested singing method is not supported")
		}

		b, err := ioutil.ReadFile(s.pubkeyPath)
		if err != nil {
			return nil, err
		}
		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
		if err != nil {
			return nil, err
		}
		return verifyKey, nil
	})
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, status.Error(codes.Unauthenticated, fmt.Sprintf("failed to verify token: %s", err.Error()))
	}
	return &pb.ValidateTokenResponse{Valid: true, User: claim.User}, nil
}

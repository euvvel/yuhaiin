package sni_spoof

import (
    "crypto/tls"
    "fmt"
    "net"
    "time"
)

type SniSpoofer struct {
    TargetIP   string
    TargetPort int
    FakeSNI    string
}

func NewSniSpoofer(targetIP string, targetPort int, fakeSNI string) *SniSpoofer {
    return &SniSpoofer{
        TargetIP:   targetIP,
        TargetPort: targetPort,
        FakeSNI:    fakeSNI,
    }
}

func (s *SniSpoofer) Spoof() error {
    targetAddr := fmt.Sprintf("%s:%d", s.TargetIP, s.TargetPort)
    
    fakeConn, err := net.DialTimeout("tcp", targetAddr, 3*time.Second)
    if err != nil {
        return fmt.Errorf("fake connection failed: %v", err)
    }
    
    fakeTLS := tls.Client(fakeConn, &tls.Config{
        ServerName:         s.FakeSNI,
        InsecureSkipVerify: true,
    })
    
    go fakeTLS.Handshake()
    time.Sleep(50 * time.Millisecond)
    fakeConn.Close()
    
    time.Sleep(100 * time.Millisecond)
    
    return nil
}

func (s *SniSpoofer) GetWhitelistedConn() (net.Conn, error) {
    targetAddr := fmt.Sprintf("%s:%d", s.TargetIP, s.TargetPort)
    
    conn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
    if err != nil {
        return nil, fmt.Errorf("real connection failed: %v", err)
    }
    
    tlsConn := tls.Client(conn, &tls.Config{
        ServerName:         s.FakeSNI,
        InsecureSkipVerify: true,
    })
    
    if err := tlsConn.Handshake(); err != nil {
        tlsConn.Close()
        return nil, fmt.Errorf("TLS handshake failed: %v", err)
    }
    
    return tlsConn, nil
}

func (s *SniSpoofer) FullSpoofAndConnect() (net.Conn, error) {
    if err := s.Spoof(); err != nil {
        return nil, err
    }
    return s.GetWhitelistedConn()
}
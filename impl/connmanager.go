package impl

import (
	"ccat/iface"
	"sync"
)

type ConnManager struct {
	ConnMap map[uint32]iface.IConn // 连接管理map,id -> Conn
	mutex   sync.RWMutex
}

func (m *ConnManager) Add(conn iface.IConn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ConnMap[conn.GetConnID()] = conn
}

func (m *ConnManager) Remove(conn iface.IConn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.ConnMap[conn.GetConnID()]; ok {
		delete(m.ConnMap, conn.GetConnID())
	}
}

func (m *ConnManager) Get(connID uint32) iface.IConn {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if conn, ok := m.ConnMap[connID]; ok {
		return conn
	}
	return nil
}

func (m *ConnManager) GetSize() uint32 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return uint32(len(m.ConnMap))
}

func (m *ConnManager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for id := range m.ConnMap {
		m.ConnMap[id].Stop()
		delete(m.ConnMap, id)
	}
}

func (m *ConnManager) RemoveAndClose(connID uint32) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if conn, ok := m.ConnMap[connID]; ok {
		conn.Stop()
		delete(m.ConnMap, connID)
	}
}

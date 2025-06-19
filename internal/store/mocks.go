package store

import "context"

func NewMockStore() Storage {
	return Storage{
		Users:         &MockUserStore{},
		Organizations: NewMockOrganizationStore(),
	}
}

type MockUserStore struct {
	usersByID    map[int64]*User
	usersByEmail map[string]*User
}

func (m *MockUserStore) Create(ctx context.Context, user *User) error {
	if m.usersByID == nil {
		m.usersByID = make(map[int64]*User)
	}
	if m.usersByEmail == nil {
		m.usersByEmail = make(map[string]*User)
	}
	m.usersByID[user.ID] = user
	if user.Email != "" {
		m.usersByEmail[user.Email] = user
	}
	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	if m.usersByID == nil {
		return nil, ErrNotFound
	}
	user, ok := m.usersByID[userID]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.usersByEmail == nil {
		return nil, ErrNotFound
	}
	user, ok := m.usersByEmail[email]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

// --- Mock Organization Store ---
type MockOrganizationStore struct {
	orgs         map[int64]*Organization
	members      map[int64]map[int64]*OrganizationMember // orgID -> userID -> member
	nextOrgID    int64
	nextMemberID int64
}

func NewMockOrganizationStore() *MockOrganizationStore {
	return &MockOrganizationStore{
		orgs:         make(map[int64]*Organization),
		members:      make(map[int64]map[int64]*OrganizationMember),
		nextOrgID:    1,
		nextMemberID: 1,
	}
}

func (m *MockOrganizationStore) Create(ctx context.Context, org *Organization) error {
	if org.ID == 0 {
		org.ID = m.nextOrgID
		m.nextOrgID++
	}
	m.orgs[org.ID] = org
	return nil
}

func (m *MockOrganizationStore) AddMember(ctx context.Context, member *OrganizationMember) error {
	if member.ID == 0 {
		member.ID = m.nextMemberID
		m.nextMemberID++
	}
	if m.members[member.OrganizationID] == nil {
		m.members[member.OrganizationID] = make(map[int64]*OrganizationMember)
	}
	m.members[member.OrganizationID][member.UserID] = member
	return nil
}

func (m *MockOrganizationStore) GetByID(ctx context.Context, orgID int64) (Organization, error) {
	org, ok := m.orgs[orgID]
	if !ok {
		return Organization{}, ErrNotFound
	}
	return *org, nil
}

func (m *MockOrganizationStore) GetMembers(ctx context.Context, orgID int64) ([]OrganizationMember, error) {
	membersMap, ok := m.members[orgID]
	if !ok {
		return nil, ErrNotFound
	}
	members := make([]OrganizationMember, 0, len(membersMap))
	for _, m := range membersMap {
		members = append(members, *m)
	}
	return members, nil
}

func (m *MockOrganizationStore) GetMember(ctx context.Context, orgID, userID int64) (OrganizationMember, error) {
	membersMap, ok := m.members[orgID]
	if !ok {
		return OrganizationMember{}, ErrNotFound
	}
	member, ok := membersMap[userID]
	if !ok {
		return OrganizationMember{}, ErrNotFound
	}
	return *member, nil
}

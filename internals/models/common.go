package models

import (
	fmt "fmt"
	"slices"
	"strconv"
	"strings"

	mpb "github.com/vapusdata-ecosystem/apis/protos/models/v1alpha1"
	encryption "github.com/vapusdata-ecosystem/vapusai/core/pkgs/encryption"
	dmutils "github.com/vapusdata-ecosystem/vapusai/core/pkgs/utils"
	types "github.com/vapusdata-ecosystem/vapusai/core/types"
)

type SupportedPackageTypes struct {
	TypeId           string `json:"typeId"`
	PackageType      string `json:"packageType"`
	PackageExtension string `json:"packageExtension"`
}

type VapusBase struct {
	CreatedAt    int64    `json:"createdAt" bun:"created_at"`
	CreatedBy    string   `json:"createdBy" bun:"created_by"`
	DeletedAt    int64    `json:"deletedAt" bun:"deleted_at,nullzero"`
	DeletedBy    string   `json:"deletedBy" bun:"deleted_by"`
	UpdatedAt    int64    `json:"updatedAt" bun:"updated_at,nullzero"`
	UpdatedBy    string   `json:"updatedBy" bun:"updated_by"`
	OwnerAccount string   `json:"ownerAccount" bun:"owner_account"`
	ID           int64    `json:"id" bun:"id,pk,autoincrement"`
	VapusID      string   `json:"vId" bun:"vapus_id,unique,notnull"`
	LastAuditID  string   `json:"lastAuditId" bun:"last_audit_id"`
	ErrorLogs    []string `json:"errorLogs" bun:"error_logs"`
	// Labels       []string `json:"labels" bun:"labels,array"`
	Organization string   `json:"Organization" bun:"Organization"`
	Status       string   `json:"status" bun:"status"`
	Editors      []string `json:"editors" bun:"editors,array"`
	Scope        string   `json:"scope" bun:"scope"`
}

func (dm *VapusBase) PreSaveVapusBase(authzClaim map[string]string) {
	if dm.CreatedBy == types.EMPTYSTR {
		dm.CreatedBy = authzClaim[encryption.ClaimUserIdKey]
	}
	if dm.CreatedAt == 0 {
		dm.CreatedAt = dmutils.GetEpochTime()
	}
	if dm.OwnerAccount == types.EMPTYSTR {
		dm.OwnerAccount = authzClaim[encryption.ClaimAccountKey]
	}
	if dm.VapusID == types.EMPTYSTR {
		dm.VapusID = dmutils.GetUUID()
	}
	dm.Organization = authzClaim[encryption.ClaimOrganizationKey]
	if dm.Scope == types.EMPTYSTR {
		dm.Scope = mpb.ResourceScope_USER_SCOPE.String()
	}
	if dm.Editors == nil {
		dm.Editors = []string{authzClaim[encryption.ClaimUserIdKey]}
	} else if !slices.Contains(dm.Editors, authzClaim[encryption.ClaimUserIdKey]) {
		dm.Editors = append(dm.Editors, authzClaim[encryption.ClaimUserIdKey])
	}
}

func (dn *VapusBase) PreDeleteVapusBase(authzClaim map[string]string) {
	if dn == nil {
		return
	}
	dn.DeletedBy = authzClaim[encryption.ClaimUserIdKey]
	dn.DeletedAt = dmutils.GetEpochTime()
	dn.Status = mpb.CommonStatus_DELETED.String()
}

func (dn *VapusBase) ConvertToPbBase() *mpb.VapusBase {
	if dn == nil {
		return &mpb.VapusBase{}
	}
	return &mpb.VapusBase{
		CreatedAt:    dn.CreatedAt,
		CreatedBy:    dn.CreatedBy,
		DeletedAt:    dn.DeletedAt,
		DeletedBy:    dn.DeletedBy,
		UpdatedAt:    dn.UpdatedAt,
		UpdatedBy:    dn.UpdatedBy,
		Organization: dn.Organization,
		Status:       dn.Status,
		Account:      dn.OwnerAccount,
		Scope:        mpb.ResourceScope(mpb.ResourceScope_value[dn.Scope]),
		Editors:      dn.Editors,
	}
}

type Comments struct {
	Comment      string `json:"comment"`
	User         string `json:"user"`
	CommentedAt  int64  `json:"commentedAt"`
	Organization string `json:"Organization"`
}

func (c *Comments) ConvertFromPb(pb *mpb.Comment) *Comments {
	if pb == nil {
		return nil
	}
	return &Comments{
		Comment:      pb.GetComment(),
		User:         pb.GetUser(),
		CommentedAt:  pb.GetCommentedAt(),
		Organization: pb.GetOrganization(),
	}
}

func (c *Comments) ConvertToPb() *mpb.Comment {
	if c == nil {
		return nil
	}
	return &mpb.Comment{
		Comment:      c.Comment,
		User:         c.User,
		CommentedAt:  c.CommentedAt,
		Organization: c.Organization,
	}
}

type JWTParams struct {
	Name                string `json:"name"`
	PublicJwtKey        string `json:"publicJwtKey"`
	PrivateJwtKey       string `json:"privateJwtKey"`
	VId                 string `json:"vId"`
	SigningAlgorithm    string `json:"signingAlgorithm"`
	IsAlreadyInSecretBs bool   `json:"isAlreadyInSecretBs"`
	Status              string `json:"status"`
	GenerateInPlatform  bool   `json:"generateInPlatform"`
}

func (j *JWTParams) Reset() {
	j = nil
}

func (j *JWTParams) GetName() string {
	if j != nil && j.Name != "" {
		return j.Name
	}
	return ""
}

func (j *JWTParams) ConvertToPb() *mpb.JWTParams {
	if j != nil {
		return &mpb.JWTParams{
			Name:                j.Name,
			PublicJwtKey:        j.PublicJwtKey,
			PrivateJwtKey:       j.PrivateJwtKey,
			VId:                 j.VId,
			SigningAlgorithm:    mpb.EncryptionAlgo(mpb.EncryptionAlgo_value[j.SigningAlgorithm]),
			IsAlreadyInSecretBs: j.IsAlreadyInSecretBs,
			Status:              j.Status,
			GenerateInPlatform:  j.GenerateInPlatform,
		}
	}
	return nil
}

func (j *JWTParams) ConvertFromPb(pb *mpb.JWTParams) *JWTParams {
	if pb == nil {
		return nil
	}
	return &JWTParams{
		Name:                pb.GetName(),
		PublicJwtKey:        pb.GetPublicJwtKey(),
		PrivateJwtKey:       pb.GetPrivateJwtKey(),
		VId:                 pb.GetVId(),
		SigningAlgorithm:    pb.GetSigningAlgorithm().String(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
		Status:              pb.GetStatus(),
		GenerateInPlatform:  pb.GetGenerateInPlatform(),
	}
}

type BackendStorages struct {
	BesType       string               `json:"besType"`
	BesOnboarding string               `json:"besOnboarding"`
	BesService    string               `json:"besService"`
	NetParams     *DataSourceNetParams `json:"netParams"`
	Status        string               `json:"status"`
}

func (b *BackendStorages) ConvertToPb() *mpb.BackendStorages {
	if b != nil {
		return &mpb.BackendStorages{
			BesType:       mpb.DataSourceType(mpb.DataSourceType_value[b.BesType]),
			BesOnboarding: mpb.BackendStorageOnboarding(mpb.BackendStorageOnboarding_value[b.BesOnboarding]),
			BesService:    mpb.DataSourceServices(mpb.DataSourceServices_value[b.BesService]),
			NetParams:     b.NetParams.ConvertToPb(),
			Status:        b.Status,
		}
	}
	return nil
}

func (b *BackendStorages) ConvertFromPb(pb *mpb.BackendStorages) *BackendStorages {
	if pb == nil {
		return nil
	}
	return &BackendStorages{
		BesType:       pb.GetBesType().String(),
		BesOnboarding: pb.GetBesOnboarding().String(),
		BesService:    pb.GetBesService().String(),
		NetParams:     (&DataSourceNetParams{}).ConvertFromPb(pb.GetNetParams()),
		Status:        pb.GetStatus(),
	}
}

type AuthnOIDC struct {
	Callback            string `json:"callback"`
	ClientId            string `json:"clientId"`
	ClientSecret        string `json:"clientSecret"`
	VId                 string `json:"vId"`
	IsAlreadyInSecretBs bool   `json:"isAlreadyInSecretBs"`
	Status              string `json:"status"`
}

func (a *AuthnOIDC) ConvertToPb() *mpb.AuthnOIDC {
	if a != nil {
		return &mpb.AuthnOIDC{
			Callback:            a.Callback,
			ClientId:            a.ClientId,
			ClientSecret:        a.ClientSecret,
			VId:                 a.VId,
			IsAlreadyInSecretBs: a.IsAlreadyInSecretBs,
			Status:              a.Status,
		}
	}
	return nil
}

func (a *AuthnOIDC) ConvertFromPb(pb *mpb.AuthnOIDC) *AuthnOIDC {
	if pb == nil {
		return nil
	}
	return &AuthnOIDC{
		Callback:            pb.GetCallback(),
		ClientId:            pb.GetClientId(),
		ClientSecret:        pb.GetClientSecret(),
		VId:                 pb.GetVId(),
		IsAlreadyInSecretBs: pb.GetIsAlreadyInSecretBs(),
		Status:              pb.GetStatus(),
	}
}

type DigestVal struct {
	Algo   string `json:"algo"`
	Digest string `json:"digest"`
}

func (d *DigestVal) GetAlgo() string {
	if d != nil && d.Algo != "" {
		return d.Algo
	}
	return ""
}

func (d *DigestVal) GetDigest() string {
	if d != nil && d.Digest != "" {
		return d.Digest
	}
	return ""
}

func (d *DigestVal) ConvertToPb() *mpb.DigestVal {
	if d != nil {
		return &mpb.DigestVal{
			Algo:   mpb.HashAlgos(mpb.HashAlgos_value[d.Algo]),
			Digest: d.GetDigest(),
		}
	}
	return nil
}

func (d *DigestVal) ConvertFromPb(pb *mpb.DigestVal) *DigestVal {
	if pb == nil {
		return nil
	}
	return &DigestVal{
		Algo:   pb.GetAlgo().String(),
		Digest: pb.GetDigest(),
	}
}

type Mapper struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *Mapper) GetKey() string {
	if m != nil && m.Key != "" {
		return m.Key
	}
	return ""
}

func (m *Mapper) GetValue() string {
	if m != nil && m.Value != "" {
		return m.Value
	}
	return ""
}

func (m *Mapper) ConvertToPb() *mpb.Mapper {
	if m != nil {
		return &mpb.Mapper{
			Key:   m.GetKey(),
			Value: m.GetValue(),
		}
	}
	return nil
}

func (m *Mapper) ConvertFromPb(pb *mpb.Mapper) *Mapper {
	if pb == nil {
		return nil
	}
	return &Mapper{
		Key:   pb.GetKey(),
		Value: pb.GetValue(),
	}
}

func MapperSliceToPb(mappers []*Mapper) []*mpb.Mapper {
	if mappers == nil {
		return nil
	}
	var list []*mpb.Mapper
	for _, m := range mappers {
		list = append(list, m.ConvertToPb())
	}
	return list
}

func MapperSliceFromPb(pbs []*mpb.Mapper) []*Mapper {
	if pbs == nil {
		return nil
	}
	var list []*Mapper
	for _, pb := range pbs {
		list = append(list, (&Mapper{}).ConvertFromPb(pb))
	}
	return list
}

type BaseIdentifier struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func (b *BaseIdentifier) GetName() string {
	if b != nil && b.Name != "" {
		return b.Name
	}
	return ""
}

func (b *BaseIdentifier) GetType() string {
	if b != nil && b.Type != "" {
		return b.Type
	}
	return ""
}

func (b *BaseIdentifier) GetIdentifier() string {
	if b != nil && b.Identifier != "" {
		return b.Identifier
	}
	return ""
}

func (b *BaseIdentifier) ConvertToPb() *mpb.BaseIdentifier {
	if b != nil {
		return &mpb.BaseIdentifier{
			Name:       b.GetName(),
			Type:       b.GetType(),
			Identifier: b.GetIdentifier(),
		}
	}
	return nil
}

func (b *BaseIdentifier) ConvertFromPb(pb *mpb.BaseIdentifier) *BaseIdentifier {
	if pb == nil {
		return nil
	}
	return &BaseIdentifier{
		Name:       pb.GetName(),
		Type:       pb.GetType(),
		Identifier: pb.GetIdentifier(),
	}
}

type FrequencyTab struct {
	Frequency         string `json:"frequency" yaml:"frequency"`
	FrequencyInterval int32  `json:"frequencyInterval" yaml:"frequencyInterval"`
}

func (c *FrequencyTab) CronvertToPb() *mpb.FrequencyTab {
	if c == nil {
		return nil
	}
	return &mpb.FrequencyTab{
		Frequency: mpb.Frequency(mpb.Frequency_value[c.Frequency]),
		Interval:  c.FrequencyInterval,
	}
}

func (c *FrequencyTab) ConvertFromPb(pb *mpb.FrequencyTab) *FrequencyTab {
	if pb == nil {
		return nil
	}
	return &FrequencyTab{
		Frequency:         pb.GetFrequency().String(),
		FrequencyInterval: pb.GetInterval(),
	}
}

type CronTab struct {
	FrequencyTab []*FrequencyTab `json:"frequencyTab" yaml:"frequencyTab"`
	Minutes      []int32         `json:"minutes" yaml:"minutes"`
	Hours        []int32         `json:"hours" yaml:"hours"`
	DaysOfMonth  []int32         `json:"daysOfMonth" yaml:"daysOfMonth"`
	Months       []int32         `json:"months" yaml:"months"`
	DaysOfWeek   []int32         `json:"daysOfWeek" yaml:"daysOfWeek"`
	CronString   string          `json:"cronString" yaml:"cronString"`
}

func (c *CronTab) CronvertToPb() *mpb.CronTab {
	if c == nil {
		return nil
	}
	fTab := make([]*mpb.FrequencyTab, 0)
	if len(c.FrequencyTab) > 0 {
		for _, f := range c.FrequencyTab {
			fTab = append(fTab, f.CronvertToPb())
		}
		return &mpb.CronTab{
			FrequencyTab: fTab,
			Minutes:      c.Minutes,
			Hours:        c.Hours,
			DaysOfMonth:  c.DaysOfMonth,
			Months:       c.Months,
			DaysOfWeek:   c.DaysOfWeek,
			CronString:   c.CronString,
		}
	}
	return &mpb.CronTab{
		FrequencyTab: []*mpb.FrequencyTab{},
		Minutes:      c.Minutes,
		Hours:        c.Hours,
		DaysOfMonth:  c.DaysOfMonth,
		Months:       c.Months,
		DaysOfWeek:   c.DaysOfWeek,
		CronString:   c.CronString,
	}
}

func (c *CronTab) ConvertFromPb(pb *mpb.CronTab) *CronTab {
	if pb == nil {
		return nil
	}
	fTab := make([]*FrequencyTab, 0)
	for _, f := range pb.GetFrequencyTab() {
		fTab = append(fTab, (&FrequencyTab{}).ConvertFromPb(f))
	}
	return &CronTab{
		FrequencyTab: fTab,
		Minutes:      pb.GetMinutes(),
		Hours:        pb.GetHours(),
		DaysOfMonth:  pb.GetDaysOfMonth(),
		Months:       pb.GetMonths(),
		DaysOfWeek:   pb.GetDaysOfWeek(),
		CronString:   pb.GetCronString(),
	}
}

type VapusSchedule struct {
	CronTab     *CronTab `json:"crontab"`
	Limit       int64    `json:"limit"`
	RunAt       int64    `json:"runAt"`
	RunNow      bool     `json:"runNow"`
	RunAfter    int64    `json:"runAfter"`
	IsRecurring bool     `json:"isRecurring"`
}

func (s *VapusSchedule) GetCrontab() *CronTab {
	if s != nil {
		return s.CronTab
	}
	return nil
}

func (s *VapusSchedule) GetLimit() int64 {
	if s != nil {
		return s.Limit
	}
	return 0
}

func (s *VapusSchedule) ConvertToPb() *mpb.VapusSchedule {
	if s != nil {
		return &mpb.VapusSchedule{
			Limit:       s.GetLimit(),
			CronTab:     s.GetCrontab().CronvertToPb(),
			RunAt:       s.RunAt,
			RunNow:      s.RunNow,
			RunAfter:    s.RunAfter,
			IsRecurring: s.IsRecurring,
		}
	}
	return nil
}

func (s *VapusSchedule) ConvertFromPb(pb *mpb.VapusSchedule) *VapusSchedule {
	if pb == nil {
		return nil
	}
	return &VapusSchedule{
		Limit:       pb.GetLimit(),
		CronTab:     (&CronTab{}).ConvertFromPb(pb.GetCronTab()),
		RunAt:       pb.GetRunAt(),
		RunNow:      pb.GetRunNow(),
		IsRecurring: pb.GetIsRecurring(),
		RunAfter:    pb.GetRunAfter(),
	}
}

func (s *VapusSchedule) ConvertCronTabToExpression() string {
	var minute, hour, day, month, weekday string
	if s.CronTab.CronString != "" {
		return s.CronTab.CronString
	}
	if s.CronTab != nil && len(s.CronTab.FrequencyTab) > 0 {
		ft := s.CronTab.FrequencyTab[0]
		switch strings.ToUpper(ft.Frequency) {
		case mpb.Frequency_MINUTELY.String():
			if ft.FrequencyInterval > 0 {
				minute = fmt.Sprintf("*/%d", ft.FrequencyInterval)
			} else {
				minute = "*"
			}
			hour = "*"
			day = "*"
			month = "*"
			weekday = "*"
		case mpb.Frequency_HOURLY.String():
			minute = "0"
			if ft.FrequencyInterval > 0 {
				hour = fmt.Sprintf("*/%d", ft.FrequencyInterval)
			} else {
				hour = "*"
			}
			day = "*"
			month = "*"
			weekday = "*"
		case mpb.Frequency_WEEKLY.String():
			minute = "0"
			hour = "0"
			day = "*"
			month = "*"
			if ft.FrequencyInterval > 0 {
				weekday = fmt.Sprintf("*/%d", ft.FrequencyInterval)
			} else {
				weekday = "*"
			}
		case mpb.Frequency_MONTHLY.String():
			minute = "0"
			hour = "0"
			if ft.FrequencyInterval > 0 {
				day = fmt.Sprintf("*/%d", ft.FrequencyInterval)
			} else {
				day = "*"
			}
			month = "*"
			weekday = "*"
		case mpb.Frequency_YEARLY.String():
			minute = "0"
			hour = "0"
			day = "1"
			if ft.FrequencyInterval > 0 {
				month = fmt.Sprintf("*/%d", ft.FrequencyInterval)
			} else {
				month = "*"
			}
			weekday = "*"
		default:
			minute = "*"
			hour = "*"
			day = "*"
			month = "*"
			weekday = "*"
		}
	} else if s.CronTab != nil {
		// Use explicit arrays from the CronTab.
		minute = intsToCronExpr(s.CronTab.Minutes, "*", 60)
		hour = intsToCronExpr(s.CronTab.Hours, "*", 24)
		day = intsToCronExpr(s.CronTab.DaysOfMonth, "*", 32)
		month = intsToCronExpr(s.CronTab.Months, "*", 13)
		weekday = intsToCronExpr(s.CronTab.DaysOfWeek, "*", 8)
	} else {
		minute, hour, day, month, weekday = "*", "*", "*", "*", "*"
	}

	return fmt.Sprintf("%s %s %s %s %s", minute, hour, day, month, weekday)
}

func (s *VapusSchedule) ParseCronExpression(expr string) (*CronTab, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression (expected 5 fields, got %d): %s", len(fields), expr)
	}

	minuteField := fields[0]
	hourField := fields[1]
	dayField := fields[2]
	monthField := fields[3]
	weekdayField := fields[4]

	ct := &CronTab{
		CronString: expr,
		// We'll leave FrequencyTab nil/empty unless you want to guess intervals
	}

	// parse each field into the corresponding slice of int32
	ct.Minutes = parseCronField(minuteField)
	ct.Hours = parseCronField(hourField)
	ct.DaysOfMonth = parseCronField(dayField)
	ct.Months = parseCronField(monthField)
	ct.DaysOfWeek = parseCronField(weekdayField)

	return ct, nil
}

// parseCronField converts a single cron field string into a slice of int32.
// e.g. "0" -> [0], "1,2,3" -> [1,2,3], "*" -> [] (meaning no specific restriction).
// If the field is something like "*/5", we do not expand it to [0,5,10,...], we simply return an empty slice
// because a direct numeric list would require further logic. Adjust as needed if you want to interpret step syntax.
func parseCronField(field string) []int32 {
	// If it's "*", we interpret it as no specific values => return empty slice
	if field == "*" {
		return nil
	}

	// If it contains "*/" or other step syntax, you can either expand it or treat it as empty
	if strings.Contains(field, "/") {
		// For a naive approach, just return nil to say "any" stepping
		return nil
	}

	// If it's a comma-separated list, parse each as an int
	parts := strings.Split(field, ",")
	var result []int32
	for _, p := range parts {
		p = strings.TrimSpace(p)
		val, err := strconv.Atoi(p)
		if err == nil {
			result = append(result, int32(val))
		} else {
			// if there's a parse error, you might skip or handle differently
		}
	}
	return result
}

// intsToCronExpr converts a slice of integers to a comma-separated string,
// or returns the default if the slice is empty.
func intsToCronExpr(arr []int32, def string, maxlimit int32) string {
	if len(arr) == 0 {
		return def
	}
	var parts []string
	for _, n := range arr {
		if n < maxlimit {
			parts = append(parts, fmt.Sprintf("%d", n))
		}
	}
	return strings.Join(parts, ",")
}

// func (s *VapusSchedule) ConvertCronTabToExpression() string {
// 	var minute, hour, day, month, weekday string

// 	if s.CronTab.EveryMinute {
// 		minute = "*"
// 	} else if s.CronTab.MinuteInterval > 0 {
// 		minute = fmt.Sprintf("*/%d", s.CronTab.MinuteInterval)
// 	} else {
// 		minute = fmt.Sprintf("%d", s.CronTab.Minute)
// 	}

// 	if s.CronTab.EveryHour {
// 		hour = "*"
// 	} else if s.CronTab.HourInterval > 0 {
// 		hour = fmt.Sprintf("*/%d", s.CronTab.HourInterval)
// 	} else {
// 		hour = fmt.Sprintf("%d", s.CronTab.Hour)
// 	}

// 	if s.CronTab.EveryDay {
// 		day = "*"
// 	} else if s.CronTab.DayInterval > 0 {
// 		day = fmt.Sprintf("*/%d", s.CronTab.DayInterval)
// 	} else {
// 		day = "*"
// 	}

// 	if s.CronTab.EveryMonth {
// 		month = "*"
// 	} else if s.CronTab.MonthInterval > 0 {
// 		month = fmt.Sprintf("*/%d", s.CronTab.MonthInterval)
// 	} else {
// 		month = "*"
// 	}

// 	if s.CronTab.EveryWeekday {
// 		weekday = "1-5"
// 	} else if len(s.CronTab.SpecificDays) > 0 {
// 		var days []string
// 		for _, d := range s.CronTab.SpecificDays {
// 			days = append(days, fmt.Sprintf("%d", d))
// 		}
// 		weekday = strings.Join(days, ",")
// 	} else {
// 		weekday = "*"
// 	}

// 	return fmt.Sprintf("%s %s %s %s %s", minute, hour, day, month, weekday)
// }

// func (s *VapusSchedule) GetCronString() (string, error) {
// 	if s == nil || s.CronTab == nil {
// 		return "", nil
// 	}

// 	minute := generateField(s.CronTab.EveryMinute, s.CronTab.MinuteInterval, "*")
// 	hour := generateField(s.CronTab.EveryHour, s.CronTab.HourInterval, "*")
// 	day := generateField(s.CronTab.EveryDay, s.CronTab.DayInterval, "*")
// 	month := generateField(s.CronTab.EveryMonth, s.CronTab.MonthInterval, "*")
// 	weekday := generateWeekdayField(s.CronTab.EveryWeekday, s.CronTab.SpecificDays)

// 	crontab := fmt.Sprintf("%s %s %s %s %s",
// 		minute, hour, day, month, weekday,
// 	)
// 	return crontab, nil
// }

// func generateField(every bool, interval int32, defaultValue string) string {
// 	if every {
// 		return defaultValue
// 	}
// 	if interval > 0 {
// 		return fmt.Sprintf("*/%d", interval)
// 	}
// 	return defaultValue
// }

// func generateWeekdayField(everyWeekday bool, specificDays []int32) string {
// 	if everyWeekday {
// 		return "1-5" // Monday to Friday
// 	}
// 	if len(specificDays) > 0 {
// 		var days []string
// 		for _, day := range specificDays {
// 			days = append(days, strconv.Itoa(int(day)))
// 		}
// 		return strings.Join(days, ",")
// 	}
// 	return "*"
// }

type QueryPrompts struct {
	Query   string   `json:"query,omitempty" yaml:"query,omitempty"`
	Schemas []string `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

func (x *QueryPrompts) ConvertToPb() *mpb.QueryPrompts {
	if x == nil {
		return nil
	}
	return &mpb.QueryPrompts{
		Query:   x.Query,
		Schemas: x.Schemas,
	}
}

func (x *QueryPrompts) ConvertFromPb(pb *mpb.QueryPrompts) *QueryPrompts {
	if pb == nil {
		return nil
	}
	return &QueryPrompts{
		Query:   pb.Query,
		Schemas: pb.Schemas,
	}
}

func GetMapperPbList(mappers []*Mapper) []*mpb.Mapper {
	if mappers == nil {
		return nil
	}
	var list []*mpb.Mapper
	for _, m := range mappers {
		list = append(list, m.ConvertToPb())
	}
	return list
}

func GetMapperObjList(pbs []*mpb.Mapper) []*Mapper {
	if pbs == nil {
		return nil
	}
	var list []*Mapper
	for _, pb := range pbs {
		list = append(list, (&Mapper{}).ConvertFromPb(pb))
	}
	return list
}

func SetVapusId(name string) string {
	return dmutils.SlugifyBase(name)
}

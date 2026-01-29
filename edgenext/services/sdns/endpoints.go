package sdns

// DNS Domain API Endpoints
const (
	EndpointDnsDomainList        = "/api/v5/sdns/domains"              // GET
	EndpointDnsDomainAdd         = "/api/v5/sdns/domains"              // POST
	EndpointDnsDomainBatchAdd    = "/api/v5/sdns/domains_batch_add"    // POST
	EndpointDnsDomainBatchDelete = "/api/v5/sdns/domains_batch_delete" // DELETE
	EndpointDnsDomainStat        = "/api/v5/sdns/domains/stat"         // GET
	EndpointDnsDomainServers     = "/api/v5/sdns/domains/servers"      // GET
)

// DNS Domain Group API Endpoints
const (
	EndpointDnsGroupList              = "/api/v5/sdns/domains/groups"                              // GET
	EndpointDnsGroupAdd               = "/api/v5/sdns/domains/groups"                              // POST
	EndpointDnsGroupUpdate            = "/api/v5/sdns/domains/groups"                              // PUT
	EndpointDnsGroupDelete            = "/api/v5/sdns/domains/groups"                              // DELETE
	EndpointDnsGroupRecordList        = "/api/v5/sdns/records/groups"                              // GET
	EndpointDnsGroupRecordRelation    = "/api/v5/sdns/records/groups_relations"                    // POST
	EndpointDnsGroupDomainList        = "/api/v5/cloud.dns.domain.group.domain.list"               // POST
	EndpointDnsGroupUndistributedList = "/api/v5/cloud.dns.domain.group.undistributed.domain.list" // POST
)

// DNS Record API Endpoints
const (
	EndpointDnsRecordTypes         = "/api/v5/sdns/types"                    // GET
	EndpointDnsRecordList          = "/api/v5/sdns/records"                  // GET
	EndpointDnsRecordAdd           = "/api/v5/sdns/records"                  // POST
	EndpointDnsRecordBatchAdd      = "/api/v5/sdns/records_batch_add"        // POST
	EndpointDnsRecordEdit          = "/api/v5/sdns/records"                  // PUT
	EndpointDnsRecordDelete        = "/api/v5/sdns/records"                  // DELETE
	EndpointDnsRecordBatchPause    = "/api/v5/sdns/records_batch_pause"      // POST
	EndpointDnsRecordBatchEnable   = "/api/v5/sdns/records_batch_enable"     // POST
	EndpointDnsRecordBatchDelete   = "/api/v5/sdns/records_batch_delete"     // POST
	EndpointDnsRecordImport        = "/api/v5/sdns/records_import"           // POST
	EndpointDnsRecordExport        = "/api/v5/sdns/records_export"           // POST
	EndpointDnsRecordLines         = "/api/v5/sdns/records/lines"            // GET
	EndpointDnsRecordGroupList     = "/api/v5/sdns/records/groups"           // GET
	EndpointDnsRecordGroupAdd      = "/api/v5/sdns/records/groups"           // POST
	EndpointDnsRecordGroupRelation = "/api/v5/sdns/records/groups_relations" // POST
	EndpointDnsRecordGroupDelete   = "/api/v5/sdns/records/groups"           // DELETE
)

// DNS Batch Task API Endpoints
const (
	EndpointDnsBatchTaskList   = "/api/v5/sdns/tasks"        // GET
	EndpointDnsBatchTaskDetail = "/api/v5/sdns/tasks/detail" // GET
)

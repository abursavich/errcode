// SPDX-License-Identifier: MIT
//
// Copyright 2025 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

// Package mysqlerr provides the ability to extract the status code from MySQL errors
// from the github.com/go-sql-driver/mysql package.
package mysqlerr

import (
	"errors"

	"bursavich.dev/errcode"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
)

var errorCoder errcode.ErrorCoder = errcode.FromFunc(ErrorCode)

// ErrorCoder return the MySQL ErrorCoder.
func ErrorCoder() errcode.ErrorCoder {
	return errorCoder
}

// SEE: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html

var mysqlCodes = map[uint16]codes.Code{
	1317: codes.Canceled, // ER_QUERY_INTERRUPTED; Query execution was interrupted

	1149: codes.InvalidArgument, // ER_SYNTAX_ERROR; You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use

	1205: codes.DeadlineExceeded, // ER_LOCK_WAIT_TIMEOUT; Lock wait timeout exceeded; try restarting transaction

	1008: codes.NotFound, // ER_DB_DROP_EXISTS; Can't drop database '%s'; database doesn't exist
	1017: codes.NotFound, // ER_FILE_NOT_FOUND; Can't find file: '%s' (errno: %d - %s)
	1031: codes.NotFound, // ER_KEY_NOT_FOUND; Can't find record in '%s'
	1049: codes.NotFound, // ER_BAD_DB_ERROR; Unknown database '%s'
	1051: codes.NotFound, // ER_BAD_TABLE_ERROR; Unknown table '%s'
	1106: codes.NotFound, // ER_UNKNOWN_PROCEDURE; Unknown procedure '%s'
	1109: codes.NotFound, // ER_UNKNOWN_TABLE; Unknown table '%s' in %s
	1133: codes.NotFound, // ER_PASSWORD_NO_MATCH; Can't find any matching row in the user table
	1146: codes.NotFound, // ER_NO_SUCH_TABLE; Table '%s.%s' doesn't exist
	1176: codes.NotFound, // ER_KEY_DOES_NOT_EXITS; Key '%s' doesn't exist in table '%s'
	1305: codes.NotFound, // ER_SP_DOES_NOT_EXIST; %s %s does not exist

	1007: codes.AlreadyExists, // ER_DB_CREATE_EXISTS; Can't create database '%s'; database exists
	1022: codes.AlreadyExists, // ER_DUP_KEY; Can't write; duplicate key in table '%s'
	1050: codes.AlreadyExists, // ER_TABLE_EXISTS_ERROR;Table '%s' already exists
	1086: codes.AlreadyExists, // ER_FILE_EXISTS_ERROR; File '%s' already exists
	1169: codes.AlreadyExists, // ER_DUP_UNIQUE; Can't write, because of unique constraint, to table '%s'
	1304: codes.AlreadyExists, // ER_SP_ALREADY_EXISTS; %s %s already exists

	1044:  codes.PermissionDenied, // ER_DBACCESS_DENIED_ERROR; Access denied for user '%s'@'%s' to database '%s'
	1045:  codes.PermissionDenied, // ER_ACCESS_DENIED_ERROR; Access denied for user '%s'@'%s' (using password: %s)
	1130:  codes.PermissionDenied, // ER_HOST_NOT_PRIVILEGED; Host '%s' is not allowed to connect to this MySQL server
	1132:  codes.PermissionDenied, // ER_PASSWORD_NOT_ALLOWED; You must have privileges to update tables in the mysql database to be able to change passwords for others
	1142:  codes.PermissionDenied, // ER_TABLEACCESS_DENIED_ERROR; %s command denied to user '%s'@'%s' for table '%s'
	1143:  codes.PermissionDenied, // ER_COLUMNACCESS_DENIED_ERROR; %s command denied to user '%s'@'%s' for column '%s' in table '%s'
	1227:  codes.PermissionDenied, // ER_SPECIFIC_ACCESS_DENIED_ERROR; SQLSTATE: Access denied; you need (at least one of) the %s privilege(s) for this operation
	1698:  codes.PermissionDenied, // ER_ACCESS_DENIED_NO_PASSWORD_ERROR; Access denied for user '%s'@'%s'
	3118:  codes.PermissionDenied, // ER_ACCOUNT_HAS_BEEN_LOCKED; Access denied for user '%s'@'%s'. Account is locked.
	3879:  codes.PermissionDenied, // ER_DB_ACCESS_DENIED; Access denied for AuthId `%s`@`%s` to database '%s
	3955:  codes.PermissionDenied, // ER_USER_ACCESS_DENIED_FOR_USER_ACCOUNT_BLOCKED_BY_PASSWORD_LOCK; Access denied for user '%s'@'%s'. Account is blocked for %s day(s) (%s day(s) remaining) due to %u consecutive failed logins.
	10926: codes.PermissionDenied, // ER_ACCESS_DENIED_ERROR_WITH_PASSWORD; Access denied for user '%s'@'%s' (using password: %s)
	10927: codes.PermissionDenied, // ER_ACCESS_DENIED_FOR_USER_ACCOUNT_LOCKED; Access denied for user '%s'@'%s'. Account is locked.
	11192: codes.PermissionDenied, // ER_FIREWALL_ACCESS_DENIED; ACCESS DENIED for '%s'. Reason: %s Statement: %s
	13525: codes.PermissionDenied, // ER_ACCESS_DENIED_FOR_USER_ACCOUNT_BLOCKED_BY_PASSWORD_LOCK; Access denied for user '%s'@'%s'. Account is blocked for %s day(s) (%s day(s) remaining) due to %u consecutive failed logins. Use FLUSH PRIVILEGES or ALTER USER to reset.

	1037: codes.ResourceExhausted, // ER_OUTOFMEMORY; Out of memory; restart server and try again (needed %d bytes)
	1038: codes.ResourceExhausted, // ER_OUT_OF_SORTMEMORY; Out of sort memory, consider increasing server sort buffer size
	1040: codes.ResourceExhausted, // ER_CON_COUNT_ERROR; Too many connections
	1041: codes.ResourceExhausted, // ER_OUT_OF_RESOURCES; Out of memory; check if mysqld or some other process uses all available memory; if not, you may have to use 'ulimit' to allow mysqld to use more memory or you can add more swap space
	1129: codes.ResourceExhausted, // ER_HOST_IS_BLOCKED; Host '%s' is blocked because of many connection errors; unblock with 'mysqladmin flush-hosts'
	1197: codes.ResourceExhausted, // ER_TRANS_CACHE_FULL; Multi-statement transaction required more than 'max_binlog_cache_size' bytes of storage; increase this mysqld variable and try again
	1203: codes.ResourceExhausted, // ER_TOO_MANY_USER_CONNECTIONS; User %s already has more than 'max_user_connections' active connections
	1206: codes.ResourceExhausted, // ER_LOCK_TABLE_FULL; The total number of locks exceeds the lock table size
	1226: codes.ResourceExhausted, // ER_USER_LIMIT_REACHED; User '%s' has exceeded the '%s' resource (current value: %ld)
	1461: codes.ResourceExhausted, // ER_MAX_PREPARED_STMT_COUNT_REACHED; Can't create more than max_prepared_stmt_count statements (current value: %lu)

	1213: codes.Aborted, // ER_LOCK_DEADLOCK; Deadlock found when trying to get lock; try restarting transaction

	1148: codes.Unimplemented, // ER_NOT_ALLOWED_COMMAND; The used command is not allowed with this MySQL version
	1178: codes.Unimplemented, // ER_CHECK_NOT_IMPLEMENTED; The storage engine for the table doesn't support %s
	1235: codes.Unimplemented, // ER_NOT_SUPPORTED_YET; This version of MySQL doesn't yet support '%s'
	1295: codes.Unimplemented, // ER_UNSUPPORTED_PS; This command is not supported in the prepared statement protocol yet

	1053: codes.Unavailable, // ER_SERVER_SHUTDOWN; Server shutdown in progress
	1077: codes.Unavailable, // ER_NORMAL_SHUTDOWN; %s: Normal shutdown
	1079: codes.Unavailable, // ER_SHUTDOWN_COMPLETE; %s: Shutdown complete
	1080: codes.Unavailable, // ER_FORCING_CLOSE; %s: Forcing close of thread %ld user: '%s'
	1194: codes.Unavailable, // ER_CRASHED_ON_USAGE; Table '%s' is marked as crashed and should be repaired
	1195: codes.Unavailable, // ER_CRASHED_ON_REPAIR; Table '%s' is marked as crashed and last (automatic?) repair failed

	1131: codes.Unauthenticated, // ER_PASSWORD_ANONYMOUS_USER; You are using MySQL as an anonymous user and anonymous users are not allowed to change passwords
}

// ErrorCode returns the gRPC code associated with the given error
// if it contains a mysql.MySQLError.
func ErrorCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if e, ok := err.(*mysql.MySQLError); ok || errors.As(err, &e) {
		if code, ok := mysqlCodes[e.Number]; ok {
			return code
		}
	}
	return codes.Unknown
}

// Objective-C API for talking to client Go package.
//   gobind -lang=objc client
//
// File is generated by gobind. Do not edit.

#ifndef __Client_H__
#define __Client_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"


@class ClientConsoleRequest;
@class ClientConsoleResponse;
@class ClientEntity;
@class ClientFileInfo;
@class ClientFileUploadRequest;
@class ClientFileUploadResponse;
@class ClientStreamSink;
@protocol ClientFireForgetPipe;
@class ClientFireForgetPipe;
@protocol ClientMetadataPush;
@class ClientMetadataPush;
@protocol ClientReqResPipe;
@class ClientReqResPipe;
@protocol ClientRequestStreamPipe;
@class ClientRequestStreamPipe;
@protocol ClientStreamResponse;
@class ClientStreamResponse;

@protocol ClientFireForgetPipe <NSObject>
- (void)fireForgetFlow:(NSData* _Nullable)input;
@end

@protocol ClientMetadataPush <NSObject>
- (void)serviceRegister:(NSData* _Nullable)input;
- (void)serviceUnRegister:(NSData* _Nullable)input;
@end

@protocol ClientReqResPipe <NSObject>
- (NSData* _Nullable)reqResFlow:(NSData* _Nullable)input error:(NSError* _Nullable* _Nullable)error;
@end

@protocol ClientRequestStreamPipe <NSObject>
- (NSData* _Nullable)requestStream:(NSData* _Nullable)input error:(NSError* _Nullable* _Nullable)error;
@end

@protocol ClientStreamResponse <NSObject>
- (void)error:(NSError* _Nullable)p0;
- (void)response:(NSData* _Nullable)p0;
@end

@interface ClientConsoleRequest : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
@property (nonatomic) NSString* _Nonnull user;
@end

@interface ClientConsoleResponse : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
@property (nonatomic) NSString* _Nonnull code;
@property (nonatomic) NSString* _Nonnull message;
@end

@interface ClientEntity : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
- (BOOL)close:(NSError* _Nullable* _Nullable)error;
- (void)console:(NSString* _Nullable)str;
- (void)fireForget:(NSString* _Nullable)name to:(NSString* _Nullable)to parameters:(NSData* _Nullable)parameters;
- (void)onMetadataPush:(id<ClientMetadataPush> _Nullable)callback;
- (BOOL)registerConsole:(NSError* _Nullable* _Nullable)error;
- (BOOL)registerConsoleListener:(NSError* _Nullable* _Nullable)error;
- (BOOL)registerFireForgetPipe:(NSString* _Nullable)name requestSchema:(NSString* _Nullable)requestSchema action:(id<ClientFireForgetPipe> _Nullable)action error:(NSError* _Nullable* _Nullable)error;
- (BOOL)registerReqResPipe:(NSString* _Nullable)name requestSchema:(NSString* _Nullable)requestSchema responseSchema:(NSString* _Nullable)responseSchema action:(id<ClientReqResPipe> _Nullable)action error:(NSError* _Nullable* _Nullable)error;
- (ClientStreamSink* _Nullable)registerReqStreamPipe:(NSString* _Nullable)name requestSchema:(NSString* _Nullable)requestSchema responseSchema:(NSString* _Nullable)responseSchema action:(id<ClientRequestStreamPipe> _Nullable)action error:(NSError* _Nullable* _Nullable)error;
- (BOOL)registerUploadFile:(NSError* _Nullable* _Nullable)error;
- (NSData* _Nullable)requestResponse:(NSString* _Nullable)name to:(NSString* _Nullable)to parameters:(NSData* _Nullable)parameters error:(NSError* _Nullable* _Nullable)error;
- (BOOL)requestStream:(NSString* _Nullable)name to:(NSString* _Nullable)to parameters:(NSData* _Nullable)parameters responseF:(id<ClientStreamResponse> _Nullable)responseF error:(NSError* _Nullable* _Nullable)error;
- (BOOL)soundOutService:(NSString* _Nullable)sever service:(NSString* _Nullable)service error:(NSError* _Nullable* _Nullable)error;
@end

@interface ClientFileInfo : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
@property (nonatomic) NSString* _Nonnull fileName;
@property (nonatomic) long fileSize;
@property (nonatomic) NSString* _Nonnull flag;
@end

@interface ClientFileUploadRequest : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
// skipped field FileUploadRequest.FileInfo with unsupported type: client.FileInfo

// skipped field FileUploadRequest.FileData with unsupported type: *[]byte

- (NSData* _Nullable)encode:(NSError* _Nullable* _Nullable)error;
@end

@interface ClientFileUploadResponse : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
@property (nonatomic) NSString* _Nonnull code;
@property (nonatomic) NSString* _Nonnull message;
// skipped field FileUploadResponse.Data with unsupported type: client.FileInfo

@end

@interface ClientStreamSink : NSObject <goSeqRefInterface> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (nonnull instancetype)init;
- (void)error:(NSError* _Nullable)err;
- (BOOL)sink:(NSData* _Nullable)data error:(NSError* _Nullable* _Nullable)error;
- (BOOL)validate:(NSData* _Nullable)data error:(NSError* _Nullable* _Nullable)error;
@end

FOUNDATION_EXPORT ClientEntity* _Nullable ClientSetup(NSString* _Nullable identify, NSString* _Nullable host, long port, id<ClientMetadataPush> _Nullable onMetadataPush, NSError* _Nullable* _Nullable error);

@class ClientFireForgetPipe;

@class ClientMetadataPush;

@class ClientReqResPipe;

@class ClientRequestStreamPipe;

@class ClientStreamResponse;

@interface ClientFireForgetPipe : NSObject <goSeqRefInterface, ClientFireForgetPipe> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)fireForgetFlow:(NSData* _Nullable)input;
@end

@interface ClientMetadataPush : NSObject <goSeqRefInterface, ClientMetadataPush> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)serviceRegister:(NSData* _Nullable)input;
- (void)serviceUnRegister:(NSData* _Nullable)input;
@end

@interface ClientReqResPipe : NSObject <goSeqRefInterface, ClientReqResPipe> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (NSData* _Nullable)reqResFlow:(NSData* _Nullable)input error:(NSError* _Nullable* _Nullable)error;
@end

@interface ClientRequestStreamPipe : NSObject <goSeqRefInterface, ClientRequestStreamPipe> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (NSData* _Nullable)requestStream:(NSData* _Nullable)input error:(NSError* _Nullable* _Nullable)error;
@end

@interface ClientStreamResponse : NSObject <goSeqRefInterface, ClientStreamResponse> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
- (void)error:(NSError* _Nullable)p0;
- (void)response:(NSData* _Nullable)p0;
@end

#endif

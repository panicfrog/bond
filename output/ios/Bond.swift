
import Foundation

public typealias UploadFileClosure = (Result<ClientFileUploadRequest, Error>) -> Void

public final class Bond {
    let port: Int
    let host: String
    let identify: String
    private var entity: ClientEntity?
    private var compeletedHandler: UploadFileHandler?
    
    public init(host: String, port: Int, identify: String) {
        self.port = port
        self.host = host
        self.identify = identify
    }
    
    public func setup() throws {
        var e: NSError?
        let client = ClientSetup(identify, host, port, nil, &e)
        if let e {
            throw e
        }
        entity = client
    }
    
    public func registerUploadService(completed: UploadFileClosure? = nil) throws {
        assert(entity != nil, "please setup first!")
        guard let entity, let completed else { return }
        compeletedHandler = UploadFileHandler(completed)
        try entity.registerUploadFile(compeletedHandler)
    }
}

enum UploadFileHandlerError: Error {
    case invalidValue
}


private class UploadFileHandler: NSObject, ClientUploadFileHandlerProtocol {
    private var closure: UploadFileClosure?
    
    convenience init(_ closure: @escaping UploadFileClosure) {
        self.init()
        self.closure = closure
    }
    
    func handleUploadFile(_ request: ClientFileUploadRequest?, error: Error?) {
        guard let closure else { return }
        if let error {
            closure(.failure(error))
            return
        }
        guard let request else {
            closure(.failure(UploadFileHandlerError.invalidValue))
            return
        }
        closure(.success(request))
    }
}

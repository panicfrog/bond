
import Foundation

public func hello() {
    debugPrint("hello world")
}

class Bond {
    let port: Int
    let host: String
    let identify: String
    
    init(port: Int, host: String, identify: String) {
        self.port = port
        self.host = host
        self.identify = identify
    }
    
    func setup() throws -> ClientEntity {
        var e: NSError? = NSError()
        let client = ClientSetup(identify, host, port, nil, &e)
        if let e {
            throw e
        }
        return client!
    }
}

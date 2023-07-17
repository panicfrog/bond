Pod::Spec.new do |s|
    s.name             = 'Bond-Swift'
    s.version          = '0.0.1'
    s.summary          = 'Bond framework for Swift.'
    s.homepage         = 'https://github.com/your-username/Bond-Swift'
    s.license          = { :type => 'MIT', :file => 'LICENSE' }
    s.author           = { 'Your Name' => 'your-email@example.com' }
    s.source           = { :git => 'https://github.com/panicfrog/Bond.git', :tag => s.version.to_s }
    s.source_files     = '**/*.{h,hpp}', '**/*.{swift}'
    s.swift_version    = '5.7'
    s.requires_arc     = true
    s.ios.deployment_target = '11.0'
    s.vendored_frameworks = 'Bond.xcframework'
  end


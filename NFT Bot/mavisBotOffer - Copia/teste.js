const ethUtil = require('ethereumjs-util');
const sigUtil = require('@metamask/eth-sig-util');

function signTypedData(privateKey, data) {
  // Parse the typed data
  const typedData = JSON.parse(data);

  // Generate the domain separator hash
  const domainHash = sigUtil.TypedDataUtils.hashStruct(
    'EIP712Domain',
    typedData.domain,
    typedData.types,
    'V4'
  );

  // Generate the message hash
  const messageHash = sigUtil.TypedDataUtils.hashStruct(
    typedData.primaryType,
    typedData.message,
    typedData.types,
    'V4'
  );

  // Combine the domain separator hash and the message hash
  const finalHash = ethUtil.keccak256(
    Buffer.concat([
      Buffer.from('1901', 'hex'),
      domainHash,
      messageHash,
    ])
  );

  // Sign the final hash
  const signature = ethUtil.ecsign(finalHash, Buffer.from(privateKey, 'hex'));

  // Return the signature in a format that includes v, r, and s
  return ethUtil.bufferToHex(Buffer.concat([
    signature.r,
    signature.s,
    Buffer.from([signature.v])
  ]));
}

var args = process.argv.slice(2);
var privateKey = args[args.length-1]
args.pop()
var data = args.join(' ');
var jsonData = JSON.parse(data)
const signature = signTypedData(privateKey, data);
console.log(signature);

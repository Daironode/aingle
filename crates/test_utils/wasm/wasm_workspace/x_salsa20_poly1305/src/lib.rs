use adk::prelude::*;

#[adk_extern]
fn x_salsa20_poly1305_encrypt(input: XSalsa20Poly1305Encrypt) -> ExternResult<XSalsa20Poly1305EncryptedData> {
    adk::prelude::x_salsa20_poly1305_encrypt(
        input.as_key_ref_ref().to_owned(),
        input.as_data_ref().to_owned(),
    )
}

#[adk_extern]
fn x_salsa20_poly1305_decrypt(input: XSalsa20Poly1305Decrypt) -> ExternResult<Option<XSalsa20Poly1305Data>> {
    adk::prelude::x_salsa20_poly1305_decrypt(
        input.as_key_ref_ref().to_owned(),
        input.as_encrypted_data_ref().to_owned()
    )
}

#[adk_extern]
fn create_x25519_keypair(_: ()) -> ExternResult<X25519PubKey> {
    adk::prelude::create_x25519_keypair()
}

#[adk_extern]
fn x_25519_x_salsa20_poly1305_encrypt(input: X25519XSalsa20Poly1305Encrypt) -> ExternResult<XSalsa20Poly1305EncryptedData> {
    adk::prelude::x_25519_x_salsa20_poly1305_encrypt(
        input.as_sender_ref().to_owned(),
        input.as_recipient_ref().to_owned(),
        input.as_data_ref().to_owned()
    )
}

#[adk_extern]
fn x_25519_x_salsa20_poly1305_decrypt(input: X25519XSalsa20Poly1305Decrypt) -> ExternResult<Option<XSalsa20Poly1305Data>> {
    adk::prelude::x_25519_x_salsa20_poly1305_decrypt(
        input.as_recipient_ref().to_owned(),
        input.as_sender_ref().to_owned(),
        input.as_encrypted_data_ref().to_owned()
    )
}

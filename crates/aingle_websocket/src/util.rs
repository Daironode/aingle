//! internal websocket utility types and code

use crate::*;

/// Implements both sides of TryFrom SerializedBytes for the passed in item.
<<<<<<< HEAD
/// See aingle_middleware_bytes::aingle_serial! macro.
=======
/// See aingle_serialized_bytes::aingle_serial! macro.
>>>>>>> master
/// This is similar, but makes use of std::io::Error for the error type.
#[macro_export]
macro_rules! try_from_serialized_bytes {
    ($s:ident) => {
<<<<<<< HEAD
        impl ::std::convert::TryFrom<$s> for ::aingle_middleware_bytes::SerializedBytes {
            type Error = ::std::io::Error;

            fn try_from(t: $s) -> ::std::io::Result<::aingle_middleware_bytes::SerializedBytes> {
                ::aingle_middleware_bytes::encode(&t)
                    .map_err(|e| ::std::io::Error::new(::std::io::ErrorKind::Other, e))
                    .map(|bytes| {
                        ::aingle_middleware_bytes::SerializedBytes::from(
                            ::aingle_middleware_bytes::UnsafeBytes::from(bytes),
=======
        impl ::std::convert::TryFrom<$s> for ::aingle_serialized_bytes::SerializedBytes {
            type Error = ::std::io::Error;

            fn try_from(t: $s) -> ::std::io::Result<::aingle_serialized_bytes::SerializedBytes> {
                ::aingle_serialized_bytes::encode(&t)
                    .map_err(|e| ::std::io::Error::new(::std::io::ErrorKind::Other, e))
                    .map(|bytes| {
                        ::aingle_serialized_bytes::SerializedBytes::from(
                            ::aingle_serialized_bytes::UnsafeBytes::from(bytes),
>>>>>>> master
                        )
                    })
            }
        }

<<<<<<< HEAD
        impl ::std::convert::TryFrom<::aingle_middleware_bytes::SerializedBytes> for $s {
            type Error = ::std::io::Error;

            fn try_from(t: ::aingle_middleware_bytes::SerializedBytes) -> ::std::io::Result<$s> {
                ::aingle_middleware_bytes::decode(t.bytes())
=======
        impl ::std::convert::TryFrom<::aingle_serialized_bytes::SerializedBytes> for $s {
            type Error = ::std::io::Error;

            fn try_from(t: ::aingle_serialized_bytes::SerializedBytes) -> ::std::io::Result<$s> {
                ::aingle_serialized_bytes::decode(t.bytes())
>>>>>>> master
                    .map_err(|e| ::std::io::Error::new(::std::io::ErrorKind::Other, e))
            }
        }
    };
}

/// not sure if we should expose this or not
/// this is the actual wire message that is sent over the websocket.
#[derive(Debug, serde::Serialize, serde::Deserialize)]
#[serde(tag = "type")]
pub(crate) enum WireMessage {
    Signal {
        #[serde(with = "serde_bytes")]
        data: Vec<u8>,
    },
    Request {
        id: String,
        #[serde(with = "serde_bytes")]
        data: Vec<u8>,
    },
    Response {
        id: String,
        #[serde(with = "serde_bytes")]
        data: Vec<u8>,
    },
}
try_from_serialized_bytes!(WireMessage);

#[cfg(test)]
pub(crate) fn init_tracing() {
    observability::test_run().unwrap();
}

/// internal socket type
pub(crate) type RawSocket = tokio_tungstenite::WebSocketStream<tokio::net::TcpStream>;

/// internal helper to convert addrs to urls
pub(crate) fn addr_to_url(a: SocketAddr, scheme: &str) -> Url2 {
    url2!("{}://{}", scheme, a)
}

/// internal helper convert urls to socket addrs for binding / connection
pub(crate) async fn url_to_addr(url: &Url2, scheme: &str) -> Result<SocketAddr> {
    if url.scheme() != scheme || url.host_str().is_none() || url.port().is_none() {
        return Err(Error::new(
            ErrorKind::InvalidInput,
            format!("got: '{}', expected: '{}://host:port'", scheme, url),
        ));
    }

    let rendered = format!("{}:{}", url.host_str().unwrap(), url.port().unwrap());

    if let Ok(mut iter) = tokio::net::lookup_host(rendered.clone()).await {
        let mut tmp = iter.next();
        let mut fallback = None;
        loop {
            if tmp.is_none() {
                break;
            }

            if tmp.as_ref().unwrap().is_ipv4() {
                return Ok(tmp.unwrap());
            }

            fallback = tmp;
            tmp = iter.next();
        }
        if let Some(addr) = fallback {
            return Ok(addr);
        }
    }

    Err(Error::new(
        ErrorKind::InvalidInput,
        format!("could not parse '{}', as 'host:port'", rendered),
    ))
}
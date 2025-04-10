from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Pagination(_message.Message):
    __slots__ = ("page", "pageSize", "totalItems", "totalPages")
    PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGESIZE_FIELD_NUMBER: _ClassVar[int]
    TOTALITEMS_FIELD_NUMBER: _ClassVar[int]
    TOTALPAGES_FIELD_NUMBER: _ClassVar[int]
    page: int
    pageSize: int
    totalItems: int
    totalPages: int
    def __init__(self, page: _Optional[int] = ..., pageSize: _Optional[int] = ..., totalItems: _Optional[int] = ..., totalPages: _Optional[int] = ...) -> None: ...

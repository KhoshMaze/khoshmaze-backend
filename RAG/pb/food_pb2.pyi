import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetAllFoodsRequest(_message.Message):
    __slots__ = ("branchID", "page", "pageSize")
    BRANCHID_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGESIZE_FIELD_NUMBER: _ClassVar[int]
    branchID: int
    page: int
    pageSize: int
    def __init__(self, branchID: _Optional[int] = ..., page: _Optional[int] = ..., pageSize: _Optional[int] = ...) -> None: ...

class Food(_message.Message):
    __slots__ = ("id", "name", "description", "type", "isAvailable", "price")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    ISAVAILABLE_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    description: str
    type: str
    isAvailable: bool
    price: float
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., description: _Optional[str] = ..., type: _Optional[str] = ..., isAvailable: bool = ..., price: _Optional[float] = ...) -> None: ...

class GetAllFoodsResponse(_message.Message):
    __slots__ = ("foods", "paginationInfo", "extra")
    class Extra(_message.Message):
        __slots__ = ("branchID",)
        BRANCHID_FIELD_NUMBER: _ClassVar[int]
        branchID: int
        def __init__(self, branchID: _Optional[int] = ...) -> None: ...
    FOODS_FIELD_NUMBER: _ClassVar[int]
    PAGINATIONINFO_FIELD_NUMBER: _ClassVar[int]
    EXTRA_FIELD_NUMBER: _ClassVar[int]
    foods: _containers.RepeatedCompositeFieldContainer[Food]
    paginationInfo: _common_pb2.Pagination
    extra: GetAllFoodsResponse.Extra
    def __init__(self, foods: _Optional[_Iterable[_Union[Food, _Mapping]]] = ..., paginationInfo: _Optional[_Union[_common_pb2.Pagination, _Mapping]] = ..., extra: _Optional[_Union[GetAllFoodsResponse.Extra, _Mapping]] = ...) -> None: ...

class CreateFoodRequest(_message.Message):
    __slots__ = ("branchID", "name", "description", "type", "price")
    BRANCHID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    branchID: int
    name: str
    description: str
    type: str
    price: float
    def __init__(self, branchID: _Optional[int] = ..., name: _Optional[str] = ..., description: _Optional[str] = ..., type: _Optional[str] = ..., price: _Optional[float] = ...) -> None: ...

package com.epfl.dedis.hbt.service.passport.mrz

open class ValidationException(msg: String) : Exception(msg)

class ChecksumException(dataType: String, data: String, computed: Int, expected: Int) :
    ValidationException("Could not validate $dataType's checksum on $data : computed $computed but expected $expected")

package com.epfl.dedis.hbt.service.passport.mrz

import android.util.Log
import com.epfl.dedis.hbt.data.Result
import com.epfl.dedis.hbt.data.Result.Error
import com.epfl.dedis.hbt.data.Result.Success
import java.util.regex.Matcher
import java.util.regex.Pattern

object MRZExtractor {

    private const val TAG: String = "Passport Validation"

    // https://en.wikipedia.org/wiki/Machine-readable_passport
    /**
     * Group 1 : Country code
     * Group 2 : Holder's name
     */
    private val LINE_1_PATTERN = Pattern.compile("P[A-Z<]([A-Z<]{3})([A-Z<]{39})")

    /**
     * Group 1 : Passport number
     * Group 2 : Passport number's checksum
     * Group 3 : Nationality
     * Group 4 : Date of birth (YYMMDD)
     * Group 5 : Date of birth checksum
     * Group 6 : Sex (M, F or < for male, female or unspecified)
     * Group 7 : Expiration date of passport (YYMMDD)
     * Group 8 : Expiration date's checksum
     * Group 9 : Personal number (may be used by the issuing country as it desires)
     * Group 10 : Personal number's checksum (may be < if all characters are <)
     * Group 11 : Checksum on Passport number, Date of birth, Expiration date and there checksums
     */
    private val LINE_2_PATTERN =
        Pattern.compile("([A-Z\\d<]{9})(\\d)([A-Z]{3})(\\d{6})(\\d)([A-Z])(\\d{6})(\\d)([A-Z\\d<]{14})([\\d<])(\\d)")

    fun match(text: String): Result<MRZInfo> {
        val matcher1 = LINE_1_PATTERN.matcher(text)
        val matcher2 = LINE_2_PATTERN.matcher(text)

        if (!matcher1.find()) return Error(NoMatchException("line1", text))
        if (!matcher2.find()) return Error(NoMatchException("line2", text))

        return try {
            Success(extractData(matcher1, matcher2))
        } catch (e: ValidationException) {
            Error(e)
        }
    }

    private fun extractData(line1: Matcher, line2: Matcher): MRZInfo {
        Log.d(
            TAG,
            "Validating passport date on lines :\n  ${line1.group()}\n  ${line2.group()}"
        )

        // Extract data adn validate them with checksums
        val (number, numberCheck) = line2.extractAndCheck("passport number", 1)
        val (dateOfBirth, birthCheck) = line2.extractAndCheck("date of birth", 4)
        val (expiration, expCheck) = line2.extractAndCheck("expiration date", 7)

        val totalData =
            number + numberCheck + dateOfBirth + birthCheck + expiration + expCheck
        val totalChecksum = line2.group(11)!!.toInt()
        validateChecksum("whole passport", totalData, totalChecksum)

        // Remove < in the pass and make sure they were at the end
        val passNumber = number.replace("<", "")
        if (!number.startsWith(passNumber))
            throw ValidationException("There were '<' in the middle of the passport number $number")

        // Extract name
        val (surname, name) = extractName(line1)
        val country = line1.group(1)!!

        return MRZInfo(
            country,
            surname,
            name,
            passNumber,
            dateOfBirth,
            expiration
        )
    }

    private fun extractName(line1: Matcher): Pair<String, String> {
        val split =
            line1.group(2)!!
                .replace('<', ' ') // Replace with whitespace
                .dropLastWhile { it.isWhitespace() } // Remove trailing spaces
                .split("  ") // split first name and last name

        if (split.isEmpty()) throw ValidationException("No name information could be retrieved")

        val surname = split[0]

        if (split.size == 1) {
            Log.d(TAG, "The holder does not have a surname")
            return "" to surname
        }

        return surname to split[1]
    }

    private fun Matcher.extractAndCheck(dataType: String, groupId: Int): Pair<String, Int> {
        val data = group(groupId)!!
        // The checksum always directly follow the extracted data
        val checksum = group(groupId + 1)!!.toInt()
        validateChecksum(dataType, data, checksum)

        return data to checksum
    }

    // https://en.wikipedia.org/wiki/Machine-readable_passport#Checksum_calculation
    private fun validateChecksum(dataType: String, data: String, expected: Int) {
        var sum = 0
        data.forEachIndexed { i, c ->
            val value = when {
                c == '<' -> 0
                c.isDigit() -> c.digitToInt()
                else -> c.minus('A') + 10
            }

            sum += value * weight(i)
        }

        sum %= 10

        if (sum != expected) throw ChecksumException(dataType, data, sum, expected)
    }

    // 2^(3-index % 3) - 1 (Basically, 0->7, 1->3, 2->1, 3->7, ...)
    private fun weight(index: Int) = (1 shl (3 - index % 3)) - 1
}
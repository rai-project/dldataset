package reader

const (
	kMagic = uint32(0xced7230a)
)

/*!
 * \brief decode the flag part of lrecord
 * \param rec the lrecord
 * \return the flag
 */
func decodeFlag(rec uint32) uint32 {
	return (rec >> uint32(29)) & uint32(7)
}

/*!
 * \brief decode the length part of lrecord
 * \param rec the lrecord
 * \return the length
 */
func decodeLength(rec uint32) uint32 {
	return rec & ((uint32(1) << uint32(29)) - 1)
}

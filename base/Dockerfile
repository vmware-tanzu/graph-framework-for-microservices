FROM gcr.io/nsx-sm/photon:4.0

ARG DEBUG
RUN \
if [ "x$DEBUG" = "xTRUE" ] || [ "x$DEBUG" = "xtrue" ] ; \
then \
  echo "DEBUG build" ; \
  tdnf -v -y install \
    jq=1.6-2.ph4 \
    logrotate=3.18.0-1.ph4 \
    python3-pip=3.9.1-6.ph4 \
    # Needed by crond to run as root. Provides /etc/pam.d entries \
    shadow=4.8.1-3.ph4 \
    cronie=1.5.5-1.ph4 \
  && pip3 install --no-cache-dir supervisor==4.2.2 setuptools==44.1.1 \
  && tdnf clean all ; \
else \
  tdnf -v -y install \
    jq=1.6-2.ph4 \
    logrotate=3.18.0-1.ph4 \
    python3-pip=3.9.1-6.ph4 \
    # Needed by crond to run as root. Provides /etc/pam.d entries \
    shadow=4.8.1-3.ph4 \
    cronie=1.5.5-1.ph4 \
  && pip3 install --no-cache-dir supervisor==4.2.2 setuptools==44.1.1 \
  && tdnf clean all ; \
fi \
# Allow root user to run crond \
&& sed -i '/#+ : root : cron crond :0 tty1 tty2 tty3 tty4 tty5 tty6/s/^#//g' /etc/security/access.conf \
# Remove unneeded cron jobs \
&& rm /etc/cron.d/*

=====================
webdevops/hhvm-apache
=====================

These image extends ``webdevops/hhvm`` with a apache daemon which is running on port 80 and 443

.. include:: include/general-supervisor.rst

Docker image tags
-----------------


.. include:: include/image-tag-hhvm.rst


Environment variables
---------------------

.. include:: include/environment-base.rst
.. include:: include/environment-base-app.rst
.. include:: include/environment-web.rst


Customization
-------------

.. include:: include/customization-apache.rst


Docker image layout
-------------------

.. include:: include/layout-apache.rst
.. include:: include/layout-hhvm.rst
